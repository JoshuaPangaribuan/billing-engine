package interactors

import (
	"context"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var _ usecases.MakePaymentUsecase = (*MakePaymentInteractor)(nil)

type (
	MakePaymentRepository interface {
		MakePayment(ctx context.Context, loanID uint64, weekNumber int64, amount string) error
		GetOutstandingString(ctx context.Context, loanID uint64) (string, error)
	}

	MakePaymentInteractorDependencies struct {
		MakePaymentRepository MakePaymentRepository
		Logger                *zap.SugaredLogger
		Validator             *validator.Validate
	}

	MakePaymentInteractor struct {
		repository MakePaymentRepository `validate:"required"`
		logger     *zap.SugaredLogger    `validate:"required"`
		validator  *validator.Validate   `validate:"required"`
	}
)

func NewMakePaymentInteractor(
	deps MakePaymentInteractorDependencies,
) *MakePaymentInteractor {
	if err := deps.Validator.Struct(deps); err != nil {
		panic(err)
	}

	return &MakePaymentInteractor{
		repository: deps.MakePaymentRepository,
		logger:     deps.Logger,
		validator:  deps.Validator,
	}
}

// Execute implements usecases.MakePaymentUsecase.
func (m *MakePaymentInteractor) Execute(ctx context.Context, input usecases.MakePaymentInput) (usecases.MakePaymentOutput, error) {
	// Validate input
	if err := m.validator.Struct(input); err != nil {
		m.logger.Errorw("invalid input", "error", err)
		return usecases.MakePaymentOutput{}, pkgerror.NewBusinessError("invalid input: " + err.Error())
	}

	// Make the payment
	err := m.repository.MakePayment(ctx, input.LoanID, input.WeekNumber, input.Amount)
	if err != nil {
		m.logger.Errorw("failed to make payment", "error", err, "loan_id", input.LoanID, "week_number", input.WeekNumber)
		return usecases.MakePaymentOutput{}, pkgerror.BusinessErrorFrom(err)
	}

	// Get updated outstanding amount
	outstanding, err := m.repository.GetOutstandingString(ctx, input.LoanID)
	if err != nil {
		m.logger.Errorw("failed to get outstanding amount", "error", err, "loan_id", input.LoanID)
		// Don't fail the payment if we can't get outstanding amount
		outstanding = "0"
	}

	message := "Payment processed successfully"
	if outstanding != "0" {
		message = "Payment processed successfully. Outstanding amount: " + outstanding
	}

	return usecases.MakePaymentOutput{
		LoanID:     input.LoanID,
		WeekNumber: input.WeekNumber,
		Amount:     input.Amount,
		Status:     "SUCCESS",
		Message:    message,
	}, nil
}
