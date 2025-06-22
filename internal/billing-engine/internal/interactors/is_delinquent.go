package interactors

import (
	"context"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"go.uber.org/zap"
)

var _ usecases.IsDelinquentUsecase = (*IsDelinquentInteractor)(nil)

type (
	IsDelinquentRepository interface {
		IsDelinquent(ctx context.Context, loanID uint64) (bool, error)
		GetInstallmentsForDelinquency(ctx context.Context, loanID uint64) ([]struct {
			WeekNumber int64
			Status     string
		}, error)
	}

	IsDelinquentInteractorDependencies struct {
		IsDelinquentRepository IsDelinquentRepository
		Logger                 *zap.SugaredLogger
	}

	IsDelinquentInteractor struct {
		repository IsDelinquentRepository `validate:"required"`
		logger     *zap.SugaredLogger     `validate:"required"`
	}
)

func NewIsDelinquentInteractor(
	deps IsDelinquentInteractorDependencies,
) *IsDelinquentInteractor {
	return &IsDelinquentInteractor{
		repository: deps.IsDelinquentRepository,
		logger:     deps.Logger,
	}
}

// Execute implements usecases.IsDelinquentUsecase.
func (i *IsDelinquentInteractor) Execute(ctx context.Context, loanID uint64) (usecases.IsDelinquentOutput, error) {
	// Check if customer is delinquent
	isDelinquent, err := i.repository.IsDelinquent(ctx, loanID)
	if err != nil {
		i.logger.Errorw("failed to check delinquency status", "error", err, "loan_id", loanID)
		return usecases.IsDelinquentOutput{}, pkgerror.BusinessErrorFrom(err)
	}

	// Get installments to provide more details
	installments, err := i.repository.GetInstallmentsForDelinquency(ctx, loanID)
	if err != nil {
		i.logger.Errorw("failed to get installments for delinquency details", "error", err, "loan_id", loanID)
		// Don't fail if we can't get details, just return basic delinquency status
		return usecases.IsDelinquentOutput{
			LoanID:       loanID,
			IsDelinquent: isDelinquent,
			Message:      getDelinquencyMessage(isDelinquent),
		}, nil
	}

	// Count missed installments and get missed weeks
	var missedWeeks []int64
	var totalMissed int64
	for _, inst := range installments {
		if inst.Status == "MISSED" {
			missedWeeks = append(missedWeeks, inst.WeekNumber)
			totalMissed++
		}
	}

	output := usecases.IsDelinquentOutput{
		LoanID:       loanID,
		IsDelinquent: isDelinquent,
		Message:      getDelinquencyMessage(isDelinquent),
		Details: struct {
			MissedWeeks []int64 `json:"missed_weeks,omitempty"`
			TotalMissed int64   `json:"total_missed"`
		}{
			MissedWeeks: missedWeeks,
			TotalMissed: totalMissed,
		},
	}

	return output, nil
}

func getDelinquencyMessage(isDelinquent bool) string {
	if isDelinquent {
		return "Customer is delinquent - has 2 or more consecutive missed payments"
	}
	return "Customer is not delinquent"
}
