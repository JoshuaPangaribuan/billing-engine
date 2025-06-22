package interactors

import (
	"context"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkguid"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var _ usecases.CreateCustomerUsecase = (*CreateCustomerInteractor)(nil)

type (
	CreateCustomerRepository interface {
		CreateCustomer(
			ctx context.Context, customer entity.Customer) (entity.Customer, error)
	}

	CreateCustomerInteractorDependencies struct {
		CustomerRepository CreateCustomerRepository
		Logger             *zap.SugaredLogger
		Validator          *validator.Validate
		SnowflakeGen       pkguid.Snowflake
	}

	CreateCustomerInteractor struct {
		repository   CreateCustomerRepository `validate:"required"`
		logger       *zap.SugaredLogger       `validate:"required"`
		validator    *validator.Validate      `validate:"required"`
		snowflakeGen pkguid.Snowflake         `validate:"required"`
	}
)

func NewCreateCustomerInteractor(
	deps CreateCustomerInteractorDependencies,
) *CreateCustomerInteractor {

	validate := validator.New()

	if err := validate.Struct(deps); err != nil {
		panic(err)
	}

	return &CreateCustomerInteractor{
		repository:   deps.CustomerRepository,
		logger:       deps.Logger,
		validator:    deps.Validator,
		snowflakeGen: deps.SnowflakeGen,
	}
}

// Execute implements usecases.CreateCustomerUsecase.
func (c *CreateCustomerInteractor) Execute(ctx context.Context, input usecases.CreateCustomerInput) (usecases.CreateCustomerOutput, error) {
	if err := c.validator.Struct(input); err != nil {
		c.logger.Errorw("invalid input", "error", err)
		return usecases.CreateCustomerOutput{}, pkgerror.ValidationErrorFrom(
			err,
		)
	}

	customer, err := c.repository.CreateCustomer(ctx, entity.Customer{
		ID:    c.snowflakeGen.Generate(),
		Name:  input.Name,
		Email: input.Email,
	})

	if err != nil {
		c.logger.Errorw("failed to create customer", "error", err)
		return usecases.CreateCustomerOutput{}, pkgerror.BusinessErrorFrom(
			err,
		)
	}

	return usecases.CreateCustomerOutput{
		ID:    customer.ID,
		Name:  customer.Name,
		Email: customer.Email,
	}, nil
}
