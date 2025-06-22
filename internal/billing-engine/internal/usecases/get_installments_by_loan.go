package usecases

import "context"

type (
	GetInstallmentsByLoanUsecase interface {
		Execute(ctx context.Context, loanID uint64) ([]GetInstallmentsOutput, error)
	}

	GetInstallmentsOutput struct {
		ID         uint64 `json:"id"`
		LoanID     uint64 `json:"loan_id"`
		WeekNumber int64  `json:"week_number"`
		DueDate    string `json:"due_date"`
		AmountDue  string `json:"amount_due"`
		Status     string `json:"status"`
	}
)
