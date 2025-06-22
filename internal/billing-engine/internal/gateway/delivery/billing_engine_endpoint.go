package delivery

import (
	"context"
	"strconv"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkghttp/v1"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type BillingEngineEndpoint struct {
	createCustomerUsecase        usecases.CreateCustomerUsecase
	getAllCustomerUsecase        usecases.GetAllCustomerUsecase
	createLoanUsecase            usecases.CreateLoanUsecase
	getInstallmentsByLoanUsecase usecases.GetInstallmentsByLoanUsecase
	makePaymentUsecase           usecases.MakePaymentUsecase
	isDelinquentUsecase          usecases.IsDelinquentUsecase
	getOutstandingUsecase        usecases.GetOutstandingUsecase

	logger    *zap.SugaredLogger
	validator *validator.Validate
}

func NewBillingEngineEndpoint(
	createCustomerUsecase usecases.CreateCustomerUsecase,
	getAllCustomerUsecase usecases.GetAllCustomerUsecase,
	createLoanUsecase usecases.CreateLoanUsecase,
	getInstallmentsByLoanUsecase usecases.GetInstallmentsByLoanUsecase,
	makePaymentUsecase usecases.MakePaymentUsecase,
	isDelinquentUsecase usecases.IsDelinquentUsecase,
	getOutstandingUsecase usecases.GetOutstandingUsecase,

	logger *zap.SugaredLogger,
	validator *validator.Validate,
) *BillingEngineEndpoint {
	return &BillingEngineEndpoint{
		createCustomerUsecase:        createCustomerUsecase,
		getAllCustomerUsecase:        getAllCustomerUsecase,
		createLoanUsecase:            createLoanUsecase,
		getInstallmentsByLoanUsecase: getInstallmentsByLoanUsecase,
		makePaymentUsecase:           makePaymentUsecase,
		isDelinquentUsecase:          isDelinquentUsecase,
		getOutstandingUsecase:        getOutstandingUsecase,

		logger:    logger,
		validator: validator,
	}
}

func (b *BillingEngineEndpoint) CreateCustomer(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	var input usecases.CreateCustomerInput
	if err := request.Decode(&input); err != nil {
		b.logger.Errorw("failed to decode request", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := b.validator.Struct(input); err != nil {
		b.logger.Errorw("failed to validate request", "error", err)

		return nil, pkgerror.ValidationErrorFrom(err)
	}

	output, err := b.createCustomerUsecase.Execute(ctx, input)
	if err != nil {
		b.logger.Errorw("failed to create customer", "error", err)

		return nil, err
	}

	return output, nil
}

func (b *BillingEngineEndpoint) GetAllCustomer(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {

	output, err := b.getAllCustomerUsecase.Execute(ctx)
	if err != nil {
		b.logger.Errorw("failed to get all customer", "error", err)

		return nil, err
	}

	return output, nil
}

func (b *BillingEngineEndpoint) CreateLoan(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	var input struct {
		CustomerID uint64 `json:"customer_id" validate:"required"`
	}
	if err := request.Decode(&input); err != nil {
		b.logger.Errorw("failed to decode request", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := b.validator.Struct(input); err != nil {
		b.logger.Errorw("failed to validate request", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	output, err := b.createLoanUsecase.Execute(ctx, input.CustomerID)
	if err != nil {
		b.logger.Errorw("failed to create loan", "error", err)
		return nil, err
	}

	return output, nil
}

func (b *BillingEngineEndpoint) GetInstallmentsByLoan(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	params := httprouter.ParamsFromContext(ctx)
	loanID := params.ByName("loan_id")

	loanIDUint, err := strconv.ParseUint(loanID, 10, 64)
	if err != nil {
		b.logger.Errorw("failed to parse loan_id", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	output, err := b.getInstallmentsByLoanUsecase.Execute(ctx, loanIDUint)
	if err != nil {
		b.logger.Errorw("failed to get installments by loan", "error", err)
		return nil, err
	}

	return output, nil
}

func (b *BillingEngineEndpoint) MakePayment(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	var input usecases.MakePaymentInput
	if err := request.Decode(&input); err != nil {
		b.logger.Errorw("failed to decode request", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	if err := b.validator.Struct(input); err != nil {
		b.logger.Errorw("failed to validate request", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	output, err := b.makePaymentUsecase.Execute(ctx, input)
	if err != nil {
		b.logger.Errorw("failed to make payment", "error", err)
		return nil, err
	}

	return output, nil
}

func (b *BillingEngineEndpoint) IsDelinquent(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	params := httprouter.ParamsFromContext(ctx)
	loanID := params.ByName("loan_id")

	loanIDUint, err := strconv.ParseUint(loanID, 10, 64)
	if err != nil {
		b.logger.Errorw("failed to parse loan_id", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	output, err := b.isDelinquentUsecase.Execute(ctx, loanIDUint)
	if err != nil {
		b.logger.Errorw("failed to check delinquency status", "error", err)
		return nil, err
	}

	return output, nil
}

func (b *BillingEngineEndpoint) GetOutstanding(
	ctx context.Context,
	request pkghttp.Request,
) (any, error) {
	params := httprouter.ParamsFromContext(ctx)
	customerID := params.ByName("customer_id")
	loanID := params.ByName("loan_id")

	customerIDUint, err := strconv.ParseUint(customerID, 10, 64)
	if err != nil {
		b.logger.Errorw("failed to parse customer_id", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	loanIDUint, err := strconv.ParseUint(loanID, 10, 64)
	if err != nil {
		b.logger.Errorw("failed to parse loan_id", "error", err)
		return nil, pkgerror.ValidationErrorFrom(err)
	}

	output, err := b.getOutstandingUsecase.Execute(ctx, customerIDUint, loanIDUint)
	if err != nil {
		b.logger.Errorw("failed to get outstanding amount", "error", err)
		return nil, err
	}

	return output, nil
}
