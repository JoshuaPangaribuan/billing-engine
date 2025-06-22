package usecases

import "context"

type (
	CreateCustomerUsecase interface {
		Execute(ctx context.Context, input CreateCustomerInput) (CreateCustomerOutput, error)
	}

	CreateCustomerInput struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	CreateCustomerOutput struct {
		ID    uint64 `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)
