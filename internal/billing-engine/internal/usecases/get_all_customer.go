package usecases

import "context"

type (
	GetAllCustomerUsecase interface {
		Execute(ctx context.Context) (GetAllCustomerOutput, error)
	}

	GetAllCustomerOutput struct {
		Customers []CustomerOutput `json:"customers"`
	}

	CustomerOutput struct {
		ID    uint64 `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)
