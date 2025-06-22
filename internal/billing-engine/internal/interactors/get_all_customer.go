package interactors

import (
	"context"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var _ usecases.GetAllCustomerUsecase = (*GetAllCustomerInteractor)(nil)

type (
	GetAllCustomerRepository interface {
		GetAllCustomer(ctx context.Context) ([]entity.Customer, error)
	}

	GetAllCustomerInteractorDependencies struct {
		GetAllCustomerRepository GetAllCustomerRepository
		Logger                   *zap.SugaredLogger
		Validator                *validator.Validate
	}

	GetAllCustomerInteractor struct {
		repository GetAllCustomerRepository `validate:"required"`
		logger     *zap.SugaredLogger       `validate:"required"`
		validator  *validator.Validate      `validate:"required"`
	}
)

func NewGetAllCustomerInteractor(
	deps GetAllCustomerInteractorDependencies,
) *GetAllCustomerInteractor {
	validate := validator.New()

	if err := validate.Struct(deps); err != nil {
		panic(err)
	}

	return &GetAllCustomerInteractor{
		repository: deps.GetAllCustomerRepository,
		logger:     deps.Logger,
		validator:  deps.Validator,
	}
}

// Execute implements usecases.GetAllCustomerUsecase.
func (g *GetAllCustomerInteractor) Execute(ctx context.Context) (usecases.GetAllCustomerOutput, error) {
	customers, err := g.repository.GetAllCustomer(ctx)
	if err != nil {
		g.logger.Errorw("failed to get all customer", "error", err)
		return usecases.GetAllCustomerOutput{}, pkgerror.BusinessErrorFrom(
			err,
		)
	}

	customersOutput := make([]usecases.CustomerOutput, len(customers))
	for i, customer := range customers {
		customersOutput[i] = usecases.CustomerOutput{
			ID:    customer.ID,
			Name:  customer.Name,
			Email: customer.Email,
		}
	}

	return usecases.GetAllCustomerOutput{
		Customers: customersOutput,
	}, nil
}
