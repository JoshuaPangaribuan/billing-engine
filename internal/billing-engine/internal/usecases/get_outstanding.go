package usecases

import "context"

type (
	GetOutstandingUsecase interface {
		Execute(ctx context.Context, customerID uint64, loanID uint64) (GetOutstandingOutput, error)
	}

	GetOutstandingOutput struct {
		CustomerID  uint64 `json:"customer_id"`
		LoanID      uint64 `json:"loan_id"`
		Outstanding struct {
			TotalAmount string `json:"total_amount"`
			TotalPaid   string `json:"total_paid"`
			TotalMissed string `json:"total_missed"`
		} `json:"outstanding"`
		Installments []struct {
			ID         uint64 `json:"id"`
			WeekNumber int64  `json:"week_number"`
			DueDate    string `json:"due_date"`
			AmountDue  string `json:"amount_due"`
			Status     string `json:"status"`
		} `json:"installments"`
	}
)
