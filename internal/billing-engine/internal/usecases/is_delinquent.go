package usecases

import "context"

type (
	IsDelinquentUsecase interface {
		Execute(ctx context.Context, loanID uint64) (IsDelinquentOutput, error)
	}

	IsDelinquentOutput struct {
		LoanID       uint64 `json:"loan_id"`
		IsDelinquent bool   `json:"is_delinquent"`
		Message      string `json:"message"`
		Details      struct {
			MissedWeeks []int64 `json:"missed_weeks,omitempty"`
			TotalMissed int64   `json:"total_missed"`
		} `json:"details,omitempty"`
	}
)
