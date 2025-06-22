package usecases

import (
	"context"
)

type (
	CreateLoanUsecase interface {
		Execute(ctx context.Context, customerID uint64) (CreateLoanOutput, error)
	}

	CreateLoanOutput struct {
		ID              uint64 `json:"id"`
		CustomerID      uint64 `json:"customer_id"`
		PrincipalAmount string `json:"principal_amount"`
		InterestRate    string `json:"interest_rate"`
		TermWeeks       int64  `json:"term_weeks"`
		StartDate       string `json:"start_date"` // format RFC3339
		Status          string `json:"status"`
	}
)
