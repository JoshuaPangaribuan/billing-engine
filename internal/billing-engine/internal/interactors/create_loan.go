package interactors

import (
	"context"
	"time"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkguid"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var _ usecases.CreateLoanUsecase = (*CreateLoanInteractor)(nil)

type (
	CreateLoanRepository interface {
		IsCustomerExist(ctx context.Context, customerID uint64) (bool, error)
		IsCustomerHasNonPaidLoan(ctx context.Context, customerID uint64) (bool, error)
		CreateLoan(ctx context.Context, loan entity.Loan) (entity.Loan, error)
		CreateInstallmentFromLoan(ctx context.Context, loan *entity.Loan) (bool, error)
	}

	CreateLoanInteractorDependencies struct {
		CreateLoanRepository CreateLoanRepository
		Logger               *zap.SugaredLogger
		SnowflakeGen         pkguid.Snowflake
	}

	CreateLoanInteractor struct {
		repository   CreateLoanRepository `validate:"required"`
		logger       *zap.SugaredLogger   `validate:"required"`
		snowflakeGen pkguid.Snowflake     `validate:"required"`
	}
)

func NewCreateLoanInteractor(
	deps CreateLoanInteractorDependencies,
) *CreateLoanInteractor {
	validate := validator.New()
	if err := validate.Struct(deps); err != nil {
		panic(err)
	}

	return &CreateLoanInteractor{
		repository:   deps.CreateLoanRepository,
		logger:       deps.Logger,
		snowflakeGen: deps.SnowflakeGen,
	}
}

// Execute implements usecases.CreateLoanUsecase.
func (c *CreateLoanInteractor) Execute(ctx context.Context, customerID uint64) (usecases.CreateLoanOutput, error) {

	isCustomerExist, err := c.repository.IsCustomerExist(ctx, customerID)
	if err != nil {
		c.logger.Error("failed to check if customer exist", zap.Error(err))
		return usecases.CreateLoanOutput{}, pkgerror.BusinessErrorFrom(
			err,
		)
	}

	if !isCustomerExist {
		return usecases.CreateLoanOutput{}, pkgerror.NewBusinessError(
			"customer not found",
		)
	}

	// check if customer has non paid loan
	// this to mitigate the case where customer has non paid loan
	// and then create a new loan, the installment will be created
	// but the loan will be paid, which is not what we want
	isCustomerHasNonPaidLoan, err := c.repository.IsCustomerHasNonPaidLoan(ctx, customerID)
	if err != nil {
		c.logger.Error("failed to check if customer has non paid loan", zap.Error(err))
		return usecases.CreateLoanOutput{}, pkgerror.BusinessErrorFrom(
			err,
		)
	}

	if isCustomerHasNonPaidLoan {
		return usecases.CreateLoanOutput{}, pkgerror.NewBusinessError(
			"customer has non paid loan",
		)
	}

	loan := entity.NewDisbursedLoan(customerID)
	loan.ID = c.snowflakeGen.Generate()
	createdLoan, err := c.repository.CreateLoan(ctx, *loan)
	if err != nil {
		c.logger.Error("failed to create loan", zap.Error(err))
		return usecases.CreateLoanOutput{}, pkgerror.BusinessErrorFrom(
			err,
		)
	}

	createInstallment, err := c.repository.CreateInstallmentFromLoan(ctx, &createdLoan)
	if err != nil {
		c.logger.Error("failed to create installment from loan", zap.Error(err))
		return usecases.CreateLoanOutput{}, pkgerror.BusinessErrorFrom(
			err,
		)
	}

	if !createInstallment {
		return usecases.CreateLoanOutput{}, pkgerror.NewBusinessError(
			"failed to create installment from loan",
		)
	}

	return usecases.CreateLoanOutput{
		ID:              createdLoan.ID,
		CustomerID:      createdLoan.CustomerID,
		PrincipalAmount: createdLoan.PrincipalAmount.String(),
		InterestRate:    createdLoan.InterestRate.String(),
		TermWeeks:       createdLoan.TermWeeks,
		StartDate:       createdLoan.StartDate.Format(time.RFC3339),
		Status:          string(createdLoan.Status),
	}, nil
}
