package billingengine

import (
	"database/sql"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/gateway/delivery"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/gateway/repository"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/interactors"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgsql"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkguid"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Exposed struct{}

type BillingEngineModuleDependencies struct {
	DB           *sql.DB
	Logger       *zap.SugaredLogger
	QueryBuilder pkgsql.GoquBuilder
	SnowflakeGen pkguid.Snowflake
	HttpRouter   *httprouter.Router
	Validator    *validator.Validate
}

func NewBillingEngineModule(
	dependencies BillingEngineModuleDependencies,
) *Exposed {

	// Billing Engine Repository
	repository := repository.NewBillingEngineRepository(
		dependencies.DB,
		dependencies.Logger,
		dependencies.QueryBuilder,
		dependencies.SnowflakeGen,
	)

	// Customer Usecases
	createCustomerInteractor := interactors.NewCreateCustomerInteractor(
		interactors.CreateCustomerInteractorDependencies{
			CustomerRepository: repository,
			Logger:             dependencies.Logger,
			Validator:          dependencies.Validator,
			SnowflakeGen:       dependencies.SnowflakeGen,
		},
	)

	getAllCustomerInteractor := interactors.NewGetAllCustomerInteractor(
		interactors.GetAllCustomerInteractorDependencies{
			GetAllCustomerRepository: repository,
			Logger:                   dependencies.Logger,
			Validator:                dependencies.Validator,
		},
	)

	// Loan Usecases
	createLoanInteractor := interactors.NewCreateLoanInteractor(
		interactors.CreateLoanInteractorDependencies{
			CreateLoanRepository: repository,
			Logger:               dependencies.Logger,
			SnowflakeGen:         dependencies.SnowflakeGen,
		},
	)

	getInstallmentsByLoanInteractor := interactors.NewGetInstallmentsByLoanInteractor(
		interactors.GetInstallmentsByLoanInteractorDependencies{
			GetInstallmentsRepository: repository,
			Logger:                    dependencies.Logger,
		},
	)

	// Billing Engine Core Usecases
	makePaymentInteractor := interactors.NewMakePaymentInteractor(
		interactors.MakePaymentInteractorDependencies{
			MakePaymentRepository: repository,
			Logger:                dependencies.Logger,
			Validator:             dependencies.Validator,
		},
	)

	isDelinquentInteractor := interactors.NewIsDelinquentInteractor(
		interactors.IsDelinquentInteractorDependencies{
			IsDelinquentRepository: repository,
			Logger:                 dependencies.Logger,
		},
	)

	getOutstandingInteractor := interactors.NewGetOutstandingInteractor(
		interactors.GetOutstandingInteractorDependencies{
			GetOutstandingRepository: repository,
			Logger:                   dependencies.Logger,
		},
	)

	// Billing Engine Endpoint
	billingEngineEndpoint := delivery.NewBillingEngineEndpoint(
		createCustomerInteractor,
		getAllCustomerInteractor,
		createLoanInteractor,
		getInstallmentsByLoanInteractor,
		makePaymentInteractor,
		isDelinquentInteractor,
		getOutstandingInteractor,
		dependencies.Logger,
		dependencies.Validator,
	)

	delivery.NewBillingEngineHTTPGateway(
		dependencies.HttpRouter,
		billingEngineEndpoint,
	)

	return &Exposed{}
}
