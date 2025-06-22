package interactors

import (
	"context"
	"errors"
	"testing"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	billingenginemocks "github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/mocks"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestGetOutstandingInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		customerID     uint64
		loanID         uint64
		setupMocks     func(*billingenginemocks.MockGetOutstandingRepository)
		expectedOutput usecases.GetOutstandingOutput
		expectedError  error
	}{
		{
			name:       "success - outstanding with installments",
			customerID: 1,
			loanID:     10,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(1)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(1), uint64(10)).Return(true, nil)
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(10)).Return("5000000", nil)
				installments := []entity.Installment{
					{ID: 1, WeekNumber: 1, DueDate: "2024-06-01", AmountDue: "100000", Status: entity.INSTALLMENT_PAID},
					{ID: 2, WeekNumber: 2, DueDate: "2024-06-08", AmountDue: "100000", Status: entity.INSTALLMENT_MISSED},
					{ID: 3, WeekNumber: 3, DueDate: "2024-06-15", AmountDue: "100000", Status: entity.INSTALLMENT_PENDING},
				}
				mockRepo.On("GetPendingInstallments", mock.Anything, uint64(10)).Return(installments, nil)
			},
			expectedOutput: func() usecases.GetOutstandingOutput {
				return usecases.GetOutstandingOutput{
					CustomerID: 1,
					LoanID:     10,
					Outstanding: struct {
						TotalAmount string `json:"total_amount"`
						TotalPaid   string `json:"total_paid"`
						TotalMissed string `json:"total_missed"`
					}{
						TotalAmount: "5000000",
						TotalPaid:   "100000.00",
						TotalMissed: "100000.00",
					},
					Installments: []struct {
						ID         uint64 `json:"id"`
						WeekNumber int64  `json:"week_number"`
						DueDate    string `json:"due_date"`
						AmountDue  string `json:"amount_due"`
						Status     string `json:"status"`
					}{
						{ID: 1, WeekNumber: 1, DueDate: "2024-06-01", AmountDue: "100000", Status: string(entity.INSTALLMENT_PAID)},
						{ID: 2, WeekNumber: 2, DueDate: "2024-06-08", AmountDue: "100000", Status: string(entity.INSTALLMENT_MISSED)},
						{ID: 3, WeekNumber: 3, DueDate: "2024-06-15", AmountDue: "100000", Status: string(entity.INSTALLMENT_PENDING)},
					},
				}
			}(),
			expectedError: nil,
		},
		{
			name:       "success - outstanding with empty installments",
			customerID: 2,
			loanID:     20,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(2)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(2), uint64(20)).Return(true, nil)
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(20)).Return("0", nil)
				mockRepo.On("GetPendingInstallments", mock.Anything, uint64(20)).Return([]entity.Installment{}, nil)
			},
			expectedOutput: func() usecases.GetOutstandingOutput {
				return usecases.GetOutstandingOutput{
					CustomerID: 2,
					LoanID:     20,
					Outstanding: struct {
						TotalAmount string `json:"total_amount"`
						TotalPaid   string `json:"total_paid"`
						TotalMissed string `json:"total_missed"`
					}{
						TotalAmount: "0",
						TotalPaid:   "0.00",
						TotalMissed: "0.00",
					},
					Installments: nil,
				}
			}(),
			expectedError: nil,
		},
		{
			name:       "error - customer not found",
			customerID: 3,
			loanID:     30,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(3)).Return(false, nil)
			},
			expectedOutput: usecases.GetOutstandingOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - loan does not belong to customer",
			customerID: 4,
			loanID:     40,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(4)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(4), uint64(40)).Return(false, nil)
			},
			expectedOutput: usecases.GetOutstandingOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on IsCustomerExist",
			customerID: 5,
			loanID:     50,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				repoErr := errors.New("db error")
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(5)).Return(false, repoErr)
			},
			expectedOutput: usecases.GetOutstandingOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on IsLoanBelongsToCustomer",
			customerID: 6,
			loanID:     60,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(6)).Return(true, nil)
				repoErr := errors.New("db error")
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(6), uint64(60)).Return(false, repoErr)
			},
			expectedOutput: usecases.GetOutstandingOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on GetOutstandingString",
			customerID: 7,
			loanID:     70,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(7)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(7), uint64(70)).Return(true, nil)
				repoErr := errors.New("db error")
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(70)).Return("", repoErr)
			},
			expectedOutput: usecases.GetOutstandingOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on GetPendingInstallments",
			customerID: 8,
			loanID:     80,
			setupMocks: func(mockRepo *billingenginemocks.MockGetOutstandingRepository) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(8)).Return(true, nil)
				mockRepo.On("IsLoanBelongsToCustomer", mock.Anything, uint64(8), uint64(80)).Return(true, nil)
				mockRepo.On("GetOutstandingString", mock.Anything, uint64(80)).Return("1000000", nil)
				repoErr := errors.New("db error")
				mockRepo.On("GetPendingInstallments", mock.Anything, uint64(80)).Return(nil, repoErr)
			},
			expectedOutput: usecases.GetOutstandingOutput{},
			expectedError:  &pkgerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := billingenginemocks.NewMockGetOutstandingRepository(t)
			logger := zap.NewNop().Sugar()

			tt.setupMocks(mockRepo)

			interactor := NewGetOutstandingInteractor(GetOutstandingInteractorDependencies{
				GetOutstandingRepository: mockRepo,
				Logger:                   logger,
			})

			output, err := interactor.Execute(context.Background(), tt.customerID, tt.loanID)

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
