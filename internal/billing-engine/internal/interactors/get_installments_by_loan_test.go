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

func TestGetInstallmentsByLoanInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		loanID         uint64
		setupMocks     func(*billingenginemocks.MockGetInstallmentsRepository)
		expectedOutput []usecases.GetInstallmentsOutput
		expectedError  error
	}{
		{
			name:   "success - get installments by loan",
			loanID: 1,
			setupMocks: func(mockRepo *billingenginemocks.MockGetInstallmentsRepository) {
				installments := []entity.Installment{
					{ID: 1, LoanID: 1, WeekNumber: 1, DueDate: "2024-06-01", AmountDue: "100000", Status: entity.INSTALLMENT_PENDING},
					{ID: 2, LoanID: 1, WeekNumber: 2, DueDate: "2024-06-08", AmountDue: "100000", Status: entity.INSTALLMENT_PAID},
				}
				mockRepo.On("GetInstallments", mock.Anything, uint64(1)).Return(installments, nil)
			},
			expectedOutput: []usecases.GetInstallmentsOutput{
				{ID: 1, LoanID: 1, WeekNumber: 1, DueDate: "2024-06-01", AmountDue: "100000", Status: string(entity.INSTALLMENT_PENDING)},
				{ID: 2, LoanID: 1, WeekNumber: 2, DueDate: "2024-06-08", AmountDue: "100000", Status: string(entity.INSTALLMENT_PAID)},
			},
			expectedError: nil,
		},
		{
			name:   "success - empty installments",
			loanID: 2,
			setupMocks: func(mockRepo *billingenginemocks.MockGetInstallmentsRepository) {
				mockRepo.On("GetInstallments", mock.Anything, uint64(2)).Return([]entity.Installment{}, nil)
			},
			expectedOutput: []usecases.GetInstallmentsOutput{},
			expectedError:  nil,
		},
		{
			name:   "error - repository error",
			loanID: 3,
			setupMocks: func(mockRepo *billingenginemocks.MockGetInstallmentsRepository) {
				repoErr := errors.New("db error")
				mockRepo.On("GetInstallments", mock.Anything, uint64(3)).Return(nil, repoErr)
			},
			expectedOutput: nil,
			expectedError:  &pkgerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := billingenginemocks.NewMockGetInstallmentsRepository(t)
			logger := zap.NewNop().Sugar()

			tt.setupMocks(mockRepo)

			interactor := NewGetInstallmentsByLoanInteractor(GetInstallmentsByLoanInteractorDependencies{
				GetInstallmentsRepository: mockRepo,
				Logger:                    logger,
			})

			output, err := interactor.Execute(context.Background(), tt.loanID)

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
