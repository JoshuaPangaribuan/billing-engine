package interactors

import (
	"context"
	"errors"
	"testing"

	billingenginemocks "github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/mocks"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestMakePaymentInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		input          usecases.MakePaymentInput
		setupMocks     func(*billingenginemocks.MockMakePaymentRepository)
		expectedOutput usecases.MakePaymentOutput
		expectedError  error
	}{
		{
			name: "success - payment processed with outstanding amount",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     1,
				WeekNumber: 5,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(100)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(100), uint64(1)).Return(true, nil)
				mockRepo.On("MakePayment", mock.Anything, uint64(1), int64(5), "100000").Return(nil)
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(1)).Return("400000", nil)
			},
			expectedOutput: usecases.MakePaymentOutput{
				CustomerID: 100,
				LoanID:     1,
				WeekNumber: 5,
				Amount:     "100000",
				Status:     "SUCCESS",
				Message:    "Payment processed successfully. Outstanding amount: 400000",
			},
			expectedError: nil,
		},
		{
			name: "success - payment processed with zero outstanding amount",
			input: usecases.MakePaymentInput{
				CustomerID: 200,
				LoanID:     2,
				WeekNumber: 10,
				Amount:     "50000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(200)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(200), uint64(2)).Return(true, nil)
				mockRepo.On("MakePayment", mock.Anything, uint64(2), int64(10), "50000").Return(nil)
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(2)).Return("0", nil)
			},
			expectedOutput: usecases.MakePaymentOutput{
				CustomerID: 200,
				LoanID:     2,
				WeekNumber: 10,
				Amount:     "50000",
				Status:     "SUCCESS",
				Message:    "Payment processed successfully",
			},
			expectedError: nil,
		},
		{
			name: "success - payment processed but outstanding amount retrieval fails (fallback)",
			input: usecases.MakePaymentInput{
				CustomerID: 300,
				LoanID:     3,
				WeekNumber: 3,
				Amount:     "75000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(300)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(300), uint64(3)).Return(true, nil)
				mockRepo.On("MakePayment", mock.Anything, uint64(3), int64(3), "75000").Return(nil)
				repoErr := errors.New("db error")
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(3)).Return("", repoErr)
			},
			expectedOutput: usecases.MakePaymentOutput{
				CustomerID: 300,
				LoanID:     3,
				WeekNumber: 3,
				Amount:     "75000",
				Status:     "SUCCESS",
				Message:    "Payment processed successfully",
			},
			expectedError: nil,
		},
		{
			name: "error - validation error (empty customer_id)",
			input: usecases.MakePaymentInput{
				CustomerID: 0,
				LoanID:     1,
				WeekNumber: 1,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - validation error (empty loan_id)",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     0,
				WeekNumber: 1,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - validation error (empty week_number)",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     1,
				WeekNumber: 0,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - validation error (empty amount)",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     1,
				WeekNumber: 1,
				Amount:     "",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - customer not found",
			input: usecases.MakePaymentInput{
				CustomerID: 999,
				LoanID:     1,
				WeekNumber: 1,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(999)).Return(false, nil)
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - loan does not belong to customer",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     999,
				WeekNumber: 1,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(100)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(100), uint64(999)).Return(false, nil)
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - repository error on IsCustomerExist",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     1,
				WeekNumber: 1,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				repoErr := errors.New("db error")
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(100)).Return(false, repoErr)
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - repository error on IsLoanBelongsToCustomer",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     1,
				WeekNumber: 1,
				Amount:     "100000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(100)).Return(true, nil)
				repoErr := errors.New("db error")
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(100), uint64(1)).Return(false, repoErr)
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "error - repository error on MakePayment",
			input: usecases.MakePaymentInput{
				CustomerID: 100,
				LoanID:     4,
				WeekNumber: 2,
				Amount:     "200000",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockMakePaymentRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(100)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(100), uint64(4)).Return(true, nil)
				repoErr := errors.New("payment failed")
				mockRepo.On("MakePayment", mock.Anything, uint64(4), int64(2), "200000").Return(repoErr)
			},
			expectedOutput: usecases.MakePaymentOutput{},
			expectedError:  &pkgerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := billingenginemocks.NewMockMakePaymentRepository(t)
			logger := zap.NewNop().Sugar()
			validator := validator.New()

			tt.setupMocks(mockRepo)

			interactor := NewMakePaymentInteractor(MakePaymentInteractorDependencies{
				MakePaymentRepository: mockRepo,
				Logger:                logger,
				Validator:             validator,
			})

			output, err := interactor.Execute(context.Background(), tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
