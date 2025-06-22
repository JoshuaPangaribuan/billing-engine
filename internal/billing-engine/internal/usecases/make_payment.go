package usecases

import "context"

type (
	MakePaymentUsecase interface {
		Execute(ctx context.Context, input MakePaymentInput) (MakePaymentOutput, error)
	}

	MakePaymentInput struct {
		CustomerID uint64 `json:"customer_id" validate:"required"`
		LoanID     uint64 `json:"loan_id" validate:"required"`
		WeekNumber int64  `json:"week_number" validate:"required"`
		Amount     string `json:"amount" validate:"required"`
	}

	MakePaymentOutput struct {
		CustomerID uint64 `json:"customer_id"`
		LoanID     uint64 `json:"loan_id"`
		WeekNumber int64  `json:"week_number"`
		Amount     string `json:"amount"`
		Status     string `json:"status"`
		Message    string `json:"message"`
	}
)
