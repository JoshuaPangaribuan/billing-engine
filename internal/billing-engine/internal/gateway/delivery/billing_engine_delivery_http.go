package delivery

import (
	"net/http"

	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkghttp/v1"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	basePath                  = "/billing-engine/api/v1"
	createCustomerPath        = "/customer"
	getAllCustomerPath        = "/customers"
	createLoanPath            = "/loan"
	getInstallmentsByLoanPath = "/loan/:loan_id/installments"
	makePaymentPath           = "/loan/payment"
	isDelinquentPath          = "/loan/:loan_id/delinquent"
	getOutstandingPath        = "/customer/:customer_id/loan/:loan_id/outstanding"
)

func NewBillingEngineHTTPGateway(
	httpRouter *httprouter.Router,
	logger *zap.SugaredLogger,
	billingEngineEndpoint *BillingEngineEndpoint,
) {
	server := pkghttp.NewServer(
		pkghttp.WithResponseEncoder(pkghttp.DefaultResponseEncoder),
		pkghttp.WithErrorResponseEncoder(pkghttp.DefaultErrorEncoder),
	)

	httpRouter.Handler(
		http.MethodPost,
		basePath+createCustomerPath,
		server.Serve(billingEngineEndpoint.CreateCustomer),
	)

	httpRouter.Handler(
		http.MethodGet,
		basePath+getAllCustomerPath,
		server.Serve(billingEngineEndpoint.GetAllCustomer),
	)

	httpRouter.Handler(
		http.MethodPost,
		basePath+createLoanPath,
		server.Serve(billingEngineEndpoint.CreateLoan),
	)

	httpRouter.Handler(
		http.MethodGet,
		basePath+getInstallmentsByLoanPath,
		server.Serve(billingEngineEndpoint.GetInstallmentsByLoan),
	)

	httpRouter.Handler(
		http.MethodPost,
		basePath+makePaymentPath,
		server.Serve(billingEngineEndpoint.MakePayment),
	)

	httpRouter.Handler(
		http.MethodGet,
		basePath+isDelinquentPath,
		server.Serve(billingEngineEndpoint.IsDelinquent),
	)

	httpRouter.Handler(
		http.MethodGet,
		basePath+getOutstandingPath,
		server.Serve(billingEngineEndpoint.GetOutstanding),
	)
}
