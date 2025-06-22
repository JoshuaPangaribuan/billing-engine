package interactors

import (
	"context"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var _ usecases.GetInstallmentsByLoanUsecase = (*GetInstallmentsByLoanInteractor)(nil)

type (
	GetInstallmentsRepository interface {
		GetInstallments(ctx context.Context, loanID uint64) ([]entity.Installment, error)
	}

	GetInstallmentsByLoanInteractorDependencies struct {
		GetInstallmentsRepository GetInstallmentsRepository
		Logger                    *zap.SugaredLogger
	}

	GetInstallmentsByLoanInteractor struct {
		repository GetInstallmentsRepository `validate:"required"`
		logger     *zap.SugaredLogger        `validate:"required"`
	}
)

func NewGetInstallmentsByLoanInteractor(
	deps GetInstallmentsByLoanInteractorDependencies,
) *GetInstallmentsByLoanInteractor {
	validate := validator.New()
	if err := validate.Struct(deps); err != nil {
		panic(err)
	}

	return &GetInstallmentsByLoanInteractor{
		repository: deps.GetInstallmentsRepository,
		logger:     deps.Logger,
	}
}

// Execute implements usecases.GetInstallmentsUsecase.
func (g *GetInstallmentsByLoanInteractor) Execute(ctx context.Context, loanID uint64) ([]usecases.GetInstallmentsOutput, error) {
	installments, err := g.repository.GetInstallments(ctx, loanID)
	if err != nil {
		g.logger.Error("failed to get installments", zap.Error(err))
		return nil, pkgerror.BusinessErrorFrom(err)
	}

	outputs := make([]usecases.GetInstallmentsOutput, len(installments))
	for i, installment := range installments {
		outputs[i] = usecases.GetInstallmentsOutput{
			ID:         installment.ID,
			LoanID:     installment.LoanID,
			WeekNumber: installment.WeekNumber,
			DueDate:    installment.DueDate,
			AmountDue:  installment.AmountDue,
			Status:     string(installment.Status),
		}
	}

	return outputs, nil
}
