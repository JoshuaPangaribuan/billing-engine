package interactors

import (
	"context"
	"strconv"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"go.uber.org/zap"
)

var _ usecases.GetOutstandingUsecase = (*GetOutstandingInteractor)(nil)

type (
	GetOutstandingRepository interface {
		GetOutstandingString(ctx context.Context, loanID uint64) (string, error)
		GetPendingInstallments(ctx context.Context, loanID uint64) ([]entity.Installment, error)
		IsCustomerExist(ctx context.Context, customerID uint64) (bool, error)
		IsLoanBelongsToCustomer(ctx context.Context, customerID uint64, loanID uint64) (bool, error)
	}

	GetOutstandingInteractorDependencies struct {
		GetOutstandingRepository GetOutstandingRepository
		Logger                   *zap.SugaredLogger
	}

	GetOutstandingInteractor struct {
		repository GetOutstandingRepository `validate:"required"`
		logger     *zap.SugaredLogger       `validate:"required"`
	}
)

func NewGetOutstandingInteractor(
	deps GetOutstandingInteractorDependencies,
) *GetOutstandingInteractor {
	return &GetOutstandingInteractor{
		repository: deps.GetOutstandingRepository,
		logger:     deps.Logger,
	}
}

// Execute implements usecases.GetOutstandingUsecase.
func (g *GetOutstandingInteractor) Execute(ctx context.Context, customerID uint64, loanID uint64) (usecases.GetOutstandingOutput, error) {
	// Check if customer exists
	isCustomerExist, err := g.repository.IsCustomerExist(ctx, customerID)
	if err != nil {
		g.logger.Errorw("failed to check if customer exists", "error", err, "customer_id", customerID)
		return usecases.GetOutstandingOutput{}, pkgerror.BusinessErrorFrom(err)
	}

	if !isCustomerExist {
		return usecases.GetOutstandingOutput{}, pkgerror.NewBusinessError("customer not found")
	}

	// Check if loan belongs to the customer
	isLoanBelongsToCustomer, err := g.repository.IsLoanBelongsToCustomer(ctx, customerID, loanID)
	if err != nil {
		g.logger.Errorw("failed to check if loan belongs to customer", "error", err, "customer_id", customerID, "loan_id", loanID)
		return usecases.GetOutstandingOutput{}, pkgerror.BusinessErrorFrom(err)
	}

	if !isLoanBelongsToCustomer {
		return usecases.GetOutstandingOutput{}, pkgerror.NewBusinessError("loan not found or does not belong to customer")
	}

	// Get outstanding amount
	outstandingAmount, err := g.repository.GetOutstandingString(ctx, loanID)
	if err != nil {
		g.logger.Errorw("failed to get outstanding amount", "error", err, "loan_id", loanID)
		return usecases.GetOutstandingOutput{}, pkgerror.BusinessErrorFrom(err)
	}

	// Get installments for detailed breakdown
	installments, err := g.repository.GetPendingInstallments(ctx, loanID)
	if err != nil {
		g.logger.Errorw("failed to get installments", "error", err, "loan_id", loanID)
		return usecases.GetOutstandingOutput{}, pkgerror.BusinessErrorFrom(err)
	}

	// Calculate totals
	var totalAmount, totalPaid, totalMissed string
	var installmentDetails []struct {
		ID         uint64 `json:"id"`
		WeekNumber int64  `json:"week_number"`
		DueDate    string `json:"due_date"`
		AmountDue  string `json:"amount_due"`
		Status     string `json:"status"`
	}

	// Calculate total paid and total missed from installments
	var totalPaidAmount, totalMissedAmount float64
	for _, inst := range installments {
		// Parse amount due to float for calculation
		amountDue, err := strconv.ParseFloat(inst.AmountDue, 64)
		if err != nil {
			g.logger.Warnw("failed to parse amount due", "error", err, "installment_id", inst.ID, "amount_due", inst.AmountDue)
			continue
		}

		// Add to appropriate total based on status
		switch inst.Status {
		case entity.INSTALLMENT_PAID:
			totalPaidAmount += amountDue
		case entity.INSTALLMENT_MISSED:
			totalMissedAmount += amountDue
		}

		installmentDetails = append(installmentDetails, struct {
			ID         uint64 `json:"id"`
			WeekNumber int64  `json:"week_number"`
			DueDate    string `json:"due_date"`
			AmountDue  string `json:"amount_due"`
			Status     string `json:"status"`
		}{
			ID:         inst.ID,
			WeekNumber: inst.WeekNumber,
			DueDate:    inst.DueDate,
			AmountDue:  inst.AmountDue,
			Status:     string(inst.Status),
		})
	}

	// Convert calculated amounts back to string format
	totalPaid = strconv.FormatFloat(totalPaidAmount, 'f', 2, 64)
	totalMissed = strconv.FormatFloat(totalMissedAmount, 'f', 2, 64)
	totalAmount = outstandingAmount // This should be the total loan amount

	output := usecases.GetOutstandingOutput{
		CustomerID: customerID,
		LoanID:     loanID,
		Outstanding: struct {
			TotalAmount string `json:"total_amount"`
			TotalPaid   string `json:"total_paid"`
			TotalMissed string `json:"total_missed"`
		}{
			TotalAmount: totalAmount,
			TotalPaid:   totalPaid,
			TotalMissed: totalMissed,
		},
		Installments: installmentDetails,
	}

	return output, nil
}
