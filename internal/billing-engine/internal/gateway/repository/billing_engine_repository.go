package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/gateway/repository/models"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgsql"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkguid"
	"github.com/doug-martin/goqu/v9"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type BillingEngineRepository struct {
	db           pkgsql.SQL
	logger       *zap.SugaredLogger
	queryBuilder pkgsql.GoquBuilder
	snowflakeGen pkguid.Snowflake

	// Tables
	customerTableName    string
	loanTableName        string
	installmentTableName string
	paymentTableName     string
}

func NewBillingEngineRepository(
	db pkgsql.SQL,
	logger *zap.SugaredLogger,
	queryBuilder pkgsql.GoquBuilder,
	snowflakeGen pkguid.Snowflake,
) *BillingEngineRepository {

	return &BillingEngineRepository{
		db:           db,
		logger:       logger,
		queryBuilder: queryBuilder,
		snowflakeGen: snowflakeGen,

		customerTableName:    "customers",
		loanTableName:        "loans",
		installmentTableName: "installments",
		paymentTableName:     "payments",
	}
}

func (b *BillingEngineRepository) CreateCustomer(
	ctx context.Context, customer entity.Customer) (entity.Customer, error) {

	createCustomer := models.Customer{
		ID:    sql.NullInt64{Int64: int64(customer.ID), Valid: true},
		Name:  sql.NullString{String: customer.Name, Valid: true},
		Email: sql.NullString{String: customer.Email, Valid: true},
	}

	query := b.queryBuilder.
		Insert(b.customerTableName).
		Cols(createCustomer.Columns()...).
		Vals(createCustomer.Values())

	sql, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)

		return entity.Customer{}, err
	}

	res, err := b.db.ExecContext(ctx, sql)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return entity.Customer{}, err
	}

	row, err := res.RowsAffected()
	if err != nil {
		b.logger.Errorw("failed to get last insert id", "error", err)

		return entity.Customer{}, err
	}

	if row == 0 {
		return entity.Customer{}, fmt.Errorf("failed to create customer")
	}

	return customer, nil
}

func (b *BillingEngineRepository) GetAllCustomer(ctx context.Context) ([]entity.Customer, error) {
	var customer models.Customer

	query := b.queryBuilder.
		Select(customer.Columns()...).
		From(b.customerTableName)

	sql, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)

		return nil, err
	}

	rows, err := b.db.QueryContext(ctx, sql)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)

		return nil, err
	}
	defer rows.Close()

	var customers []entity.Customer
	for rows.Next() {
		err := rows.Scan(customer.Values()...)
		if err != nil {
			b.logger.Errorw("failed to scan row", "error", err)

			return nil, err
		}

		customers = append(customers, entity.Customer{
			ID:    uint64(customer.ID.Int64),
			Name:  customer.Name.String,
			Email: customer.Email.String,
		})
	}

	return customers, nil
}

// Loan Usecases
func (b *BillingEngineRepository) IsCustomerExist(ctx context.Context, customerID uint64) (bool, error) {
	var customer models.Customer

	query := b.queryBuilder.
		Select("id").
		From(b.customerTableName).
		Where(goqu.Ex{"id": customerID})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return false, err
	}

	row := b.db.QueryRowContext(ctx, sqlQuery)
	err = row.Scan(&customer.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		b.logger.Errorw("failed to scan row", "error", err)
		return false, err
	}

	return true, nil
}

func (b *BillingEngineRepository) IsCustomerHasNonPaidLoan(ctx context.Context, customerID uint64) (bool, error) {
	var loan models.Loan

	query := b.queryBuilder.
		Select("id").
		From(b.loanTableName).
		Where(goqu.Ex{"customer_id": customerID}).
		Where(goqu.Ex{"status": "DISBURSED"})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return false, err
	}

	row := b.db.QueryRowContext(ctx, sqlQuery)
	err = row.Scan(&loan.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		b.logger.Errorw("failed to scan row", "error", err)
		return false, err
	}

	return true, nil
}

func (b *BillingEngineRepository) IsLoanBelongsToCustomer(ctx context.Context, customerID uint64, loanID uint64) (bool, error) {
	var loan models.Loan

	query := b.queryBuilder.
		Select("id").
		From(b.loanTableName).
		Where(goqu.Ex{"id": loanID}).
		Where(goqu.Ex{"customer_id": customerID})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return false, err
	}

	row := b.db.QueryRowContext(ctx, sqlQuery)
	err = row.Scan(&loan.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		b.logger.Errorw("failed to scan row", "error", err)
		return false, err
	}

	return true, nil
}

func (b *BillingEngineRepository) CreateLoan(ctx context.Context, loan entity.Loan) (entity.Loan, error) {
	createLoan := models.Loan{
		ID:              sql.NullInt64{Int64: int64(loan.ID), Valid: true},
		CustomerID:      sql.NullInt64{Int64: int64(loan.CustomerID), Valid: true},
		PrincipalAmount: loan.PrincipalAmount,
		InterestRate:    loan.InterestRate,
		TermWeeks:       sql.NullInt64{Int64: loan.TermWeeks, Valid: true},
		StartDate:       sql.NullTime{Time: loan.StartDate, Valid: true},
		Status:          sql.NullString{String: string(loan.Status), Valid: true},
	}

	query := b.queryBuilder.
		Insert(b.loanTableName).
		Cols(createLoan.Columns()...).
		Vals(createLoan.Values())

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return entity.Loan{}, err
	}

	res, err := b.db.ExecContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return entity.Loan{}, err
	}

	row, err := res.RowsAffected()
	if err != nil {
		b.logger.Errorw("failed to get rows affected", "error", err)
		return entity.Loan{}, err
	}

	if row == 0 {
		return entity.Loan{}, fmt.Errorf("failed to create loan")
	}

	return loan, nil
}

func (b *BillingEngineRepository) CreateInstallmentFromLoan(ctx context.Context, loan *entity.Loan) (bool, error) {
	// Calculate weekly payment amount (principal + interest) / number of weeks
	totalAmount := loan.PrincipalAmount.Mul(decimal.NewFromFloat(1.1)) // 10% interest
	weeklyAmount := totalAmount.Div(decimal.NewFromInt(loan.TermWeeks))

	// Create installments for each week
	for week := int64(1); week <= loan.TermWeeks; week++ {
		dueDate := loan.StartDate.AddDate(0, 0, int(week*7)) // Add weeks

		installment := models.Installment{
			ID:         sql.NullInt64{Int64: int64(b.snowflakeGen.Generate()), Valid: true},
			LoanID:     sql.NullInt64{Int64: int64(loan.ID), Valid: true},
			WeekNumber: sql.NullInt64{Int64: week, Valid: true},
			DueDate:    sql.NullString{String: dueDate.Format("2006-01-02"), Valid: true},
			AmountDue:  sql.NullString{String: weeklyAmount.String(), Valid: true},
			Status:     sql.NullString{String: "PENDING", Valid: true},
		}

		query := b.queryBuilder.
			Insert(b.installmentTableName).
			Cols(installment.Columns()...).
			Vals(installment.Values())

		sqlQuery, _, err := query.ToSQL()
		if err != nil {
			b.logger.Errorw("failed to build query", "error", err)
			return false, err
		}

		_, err = b.db.ExecContext(ctx, sqlQuery)
		if err != nil {
			b.logger.Errorw("failed to execute query", "error", err)
			return false, err
		}
	}

	return true, nil
}

// Installment Usecases
func (b *BillingEngineRepository) GetInstallments(ctx context.Context, loanID uint64) ([]entity.Installment, error) {
	var installment models.Installment

	query := b.queryBuilder.
		Select(installment.Columns()...).
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Order(goqu.C("week_number").Asc())

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return nil, err
	}

	rows, err := b.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var installments []entity.Installment
	for rows.Next() {
		err := rows.Scan(installment.Values()...)
		if err != nil {
			b.logger.Errorw("failed to scan row", "error", err)
			return nil, err
		}

		installments = append(installments, entity.Installment{
			ID:         uint64(installment.ID.Int64),
			LoanID:     uint64(installment.LoanID.Int64),
			WeekNumber: installment.WeekNumber.Int64,
			DueDate:    installment.DueDate.String,
			AmountDue:  installment.AmountDue.String,
			Status:     entity.InstallmentStatus(installment.Status.String),
		})
	}

	return installments, nil
}

// Additional methods for billing engine functionality
func (b *BillingEngineRepository) GetOutstanding(ctx context.Context, loanID uint64) (decimal.Decimal, error) {
	var installment models.Installment

	query := b.queryBuilder.
		Select("amount_due").
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Where(goqu.Ex{"status": []string{"PENDING", "MISSED"}})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return decimal.Zero, err
	}

	rows, err := b.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return decimal.Zero, err
	}
	defer rows.Close()

	totalOutstanding := decimal.Zero
	for rows.Next() {
		err := rows.Scan(&installment.AmountDue)
		if err != nil {
			b.logger.Errorw("failed to scan row", "error", err)
			return decimal.Zero, err
		}

		amount, err := decimal.NewFromString(installment.AmountDue.String)
		if err != nil {
			b.logger.Errorw("failed to parse amount", "error", err)
			return decimal.Zero, err
		}

		totalOutstanding = totalOutstanding.Add(amount)
	}

	return totalOutstanding, nil
}

func (b *BillingEngineRepository) IsDelinquent(ctx context.Context, loanID uint64) (bool, error) {
	var installment models.Installment

	// Check for 2 consecutive missed payments
	query := b.queryBuilder.
		Select("week_number").
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Where(goqu.Ex{"status": "MISSED"}).
		Order(goqu.C("week_number").Desc()).
		Limit(2)

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return false, err
	}

	rows, err := b.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return false, err
	}
	defer rows.Close()

	var missedWeeks []int64
	for rows.Next() {
		err := rows.Scan(&installment.WeekNumber)
		if err != nil {
			b.logger.Errorw("failed to scan row", "error", err)
			return false, err
		}
		missedWeeks = append(missedWeeks, installment.WeekNumber.Int64)
	}

	// Check if there are 2 consecutive missed payments
	if len(missedWeeks) >= 2 {
		// Check if the two most recent missed payments are consecutive
		if missedWeeks[0] == missedWeeks[1]+1 {
			return true, nil
		}
	}

	return false, nil
}

func (b *BillingEngineRepository) MakePayment(ctx context.Context, loanID uint64, weekNumber int64, amount string) error {
	// First, find the installment for the specified week
	var installment models.Installment

	query := b.queryBuilder.
		Select(installment.Columns()...).
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Where(goqu.Ex{"week_number": weekNumber})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return err
	}

	row := b.db.QueryRowContext(ctx, sqlQuery)
	err = row.Scan(installment.Values()...)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("installment not found for loan %d week %d", loanID, weekNumber)
		}
		b.logger.Errorw("failed to scan row", "error", err)
		return err
	}

	// Check if the payment amount matches the amount due
	if installment.AmountDue.String != amount {
		return fmt.Errorf("payment amount %s does not match amount due %s", amount, installment.AmountDue.String)
	}

	// Check if installment is already paid
	if installment.Status.String == "PAID" {
		return fmt.Errorf("installment for loan %d week %d is already paid", loanID, weekNumber)
	}

	// Create payment record
	payment := models.Payment{
		ID:            sql.NullInt64{Int64: int64(b.snowflakeGen.Generate()), Valid: true},
		InstallmentID: sql.NullInt64{Int64: installment.ID.Int64, Valid: true},
		PaidAt:        sql.NullTime{Time: time.Now(), Valid: true},
		AmountPaid:    sql.NullString{String: amount, Valid: true},
	}

	paymentQuery := b.queryBuilder.
		Insert(b.paymentTableName).
		Cols(payment.Columns()...).
		Vals(payment.Values())

	paymentSQL, _, err := paymentQuery.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build payment query", "error", err)
		return err
	}

	_, err = b.db.ExecContext(ctx, paymentSQL)
	if err != nil {
		b.logger.Errorw("failed to execute payment query", "error", err)
		return err
	}

	// Update installment status to PAID
	updateQuery := b.queryBuilder.
		Update(b.installmentTableName).
		Set(goqu.Record{"status": "PAID"}).
		Where(goqu.Ex{"id": installment.ID.Int64})

	updateSQL, _, err := updateQuery.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build update query", "error", err)
		return err
	}

	_, err = b.db.ExecContext(ctx, updateSQL)
	if err != nil {
		b.logger.Errorw("failed to execute update query", "error", err)
		return err
	}

	// Check if all installments are paid to update loan status
	allPaidQuery := b.queryBuilder.
		Select(goqu.COUNT("*")).
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Where(goqu.Ex{"status": goqu.Op{"neq": "PAID"}})

	allPaidSQL, _, err := allPaidQuery.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build all paid query", "error", err)
		return err
	}

	var unpaidCount int64
	row = b.db.QueryRowContext(ctx, allPaidSQL)
	err = row.Scan(&unpaidCount)
	if err != nil {
		b.logger.Errorw("failed to scan unpaid count", "error", err)
		return err
	}

	// If no unpaid installments, mark loan as PAID
	if unpaidCount == 0 {
		loanUpdateQuery := b.queryBuilder.
			Update(b.loanTableName).
			Set(goqu.Record{"status": "PAID"}).
			Where(goqu.Ex{"id": loanID})

		loanUpdateSQL, _, err := loanUpdateQuery.ToSQL()
		if err != nil {
			b.logger.Errorw("failed to build loan update query", "error", err)
			return err
		}

		_, err = b.db.ExecContext(ctx, loanUpdateSQL)
		if err != nil {
			b.logger.Errorw("failed to execute loan update query", "error", err)
			return err
		}
	}

	return nil
}

func (b *BillingEngineRepository) UpdateMissedInstallments(ctx context.Context, loanID uint64) error {
	// Update installments that are past due date and still pending
	query := b.queryBuilder.
		Update(b.installmentTableName).
		Set(goqu.Record{"status": "MISSED"}).
		Where(goqu.Ex{"loan_id": loanID}).
		Where(goqu.Ex{"status": "PENDING"}).
		Where(goqu.Ex{"due_date": goqu.Op{"lt": time.Now().Format("2006-01-02")}})

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return err
	}

	_, err = b.db.ExecContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return err
	}

	return nil
}

// Additional methods for interactors
func (b *BillingEngineRepository) GetOutstandingString(ctx context.Context, loanID uint64) (string, error) {
	outstanding, err := b.GetOutstanding(ctx, loanID)
	if err != nil {
		return "", err
	}
	return outstanding.String(), nil
}

func (b *BillingEngineRepository) GetInstallmentsForDelinquency(ctx context.Context, loanID uint64) ([]struct {
	WeekNumber int64
	Status     string
}, error) {
	var installment models.Installment

	query := b.queryBuilder.
		Select("week_number", "status").
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Order(goqu.C("week_number").Asc())

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return nil, err
	}

	rows, err := b.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var installments []struct {
		WeekNumber int64
		Status     string
	}
	for rows.Next() {
		err := rows.Scan(&installment.WeekNumber, &installment.Status)
		if err != nil {
			b.logger.Errorw("failed to scan row", "error", err)
			return nil, err
		}

		installments = append(installments, struct {
			WeekNumber int64
			Status     string
		}{
			WeekNumber: installment.WeekNumber.Int64,
			Status:     installment.Status.String,
		})
	}

	return installments, nil
}

func (b *BillingEngineRepository) GetAllInstallments(ctx context.Context, loanID uint64) ([]entity.Installment, error) {
	var installment models.Installment

	query := b.queryBuilder.
		Select(installment.Columns()...).
		From(b.installmentTableName).
		Where(goqu.Ex{"loan_id": loanID}).
		Order(goqu.C("week_number").Asc())

	sqlQuery, _, err := query.ToSQL()
	if err != nil {
		b.logger.Errorw("failed to build query", "error", err)
		return nil, err
	}

	rows, err := b.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		b.logger.Errorw("failed to execute query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var installments []entity.Installment
	for rows.Next() {
		err := rows.Scan(installment.Values()...)
		if err != nil {
			b.logger.Errorw("failed to scan row", "error", err)
			return nil, err
		}

		installments = append(installments, entity.Installment{
			ID:         uint64(installment.ID.Int64),
			LoanID:     uint64(installment.LoanID.Int64),
			WeekNumber: installment.WeekNumber.Int64,
			DueDate:    installment.DueDate.String,
			AmountDue:  installment.AmountDue.String,
			Status:     entity.InstallmentStatus(installment.Status.String),
		})
	}

	return installments, nil
}
