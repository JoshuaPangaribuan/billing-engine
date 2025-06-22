package interactors

import (
	"context"
	"errors"
	"testing"

	billingenginemocks "github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/mocks"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestIsDelinquentInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		loanID         uint64
		setupMocks     func(*billingenginemocks.MockIsDelinquentRepository)
		expectedOutput usecases.IsDelinquentOutput
		expectedError  error
	}{
		{
			name:   "success - customer is delinquent with details",
			loanID: 1,
			setupMocks: func(mockRepo *billingenginemocks.MockIsDelinquentRepository) {
				mockRepo.On("IsDelinquent", mock.Anything, uint64(1)).Return(true, nil)
				installments := []struct {
					WeekNumber int64
					Status     string
				}{
					{WeekNumber: 1, Status: "PAID"},
					{WeekNumber: 2, Status: "MISSED"},
					{WeekNumber: 3, Status: "MISSED"},
					{WeekNumber: 4, Status: "PENDING"},
				}
				mockRepo.On("GetInstallmentsForDelinquency", mock.Anything, uint64(1)).Return(installments, nil)
			},
			expectedOutput: usecases.IsDelinquentOutput{
				LoanID:       1,
				IsDelinquent: true,
				Message:      "Customer is delinquent - has 2 or more consecutive missed payments",
				Details: struct {
					MissedWeeks []int64 `json:"missed_weeks,omitempty"`
					TotalMissed int64   `json:"total_missed"`
				}{
					MissedWeeks: []int64{2, 3},
					TotalMissed: 2,
				},
			},
			expectedError: nil,
		},
		{
			name:   "success - customer is not delinquent",
			loanID: 2,
			setupMocks: func(mockRepo *billingenginemocks.MockIsDelinquentRepository) {
				mockRepo.On("IsDelinquent", mock.Anything, uint64(2)).Return(false, nil)
				installments := []struct {
					WeekNumber int64
					Status     string
				}{
					{WeekNumber: 1, Status: "PAID"},
					{WeekNumber: 2, Status: "PAID"},
					{WeekNumber: 3, Status: "PENDING"},
				}
				mockRepo.On("GetInstallmentsForDelinquency", mock.Anything, uint64(2)).Return(installments, nil)
			},
			expectedOutput: usecases.IsDelinquentOutput{
				LoanID:       2,
				IsDelinquent: false,
				Message:      "Customer is not delinquent",
				Details: struct {
					MissedWeeks []int64 `json:"missed_weeks,omitempty"`
					TotalMissed int64   `json:"total_missed"`
				}{
					MissedWeeks: nil,
					TotalMissed: 0,
				},
			},
			expectedError: nil,
		},
		{
			name:   "success - customer is delinquent but details unavailable (fallback)",
			loanID: 3,
			setupMocks: func(mockRepo *billingenginemocks.MockIsDelinquentRepository) {
				mockRepo.On("IsDelinquent", mock.Anything, uint64(3)).Return(true, nil)
				repoErr := errors.New("db error")
				mockRepo.On("GetInstallmentsForDelinquency", mock.Anything, uint64(3)).Return(nil, repoErr)
			},
			expectedOutput: usecases.IsDelinquentOutput{
				LoanID:       3,
				IsDelinquent: true,
				Message:      "Customer is delinquent - has 2 or more consecutive missed payments",
			},
			expectedError: nil,
		},
		{
			name:   "success - customer is not delinquent but details unavailable (fallback)",
			loanID: 4,
			setupMocks: func(mockRepo *billingenginemocks.MockIsDelinquentRepository) {
				mockRepo.On("IsDelinquent", mock.Anything, uint64(4)).Return(false, nil)
				repoErr := errors.New("db error")
				mockRepo.On("GetInstallmentsForDelinquency", mock.Anything, uint64(4)).Return(nil, repoErr)
			},
			expectedOutput: usecases.IsDelinquentOutput{
				LoanID:       4,
				IsDelinquent: false,
				Message:      "Customer is not delinquent",
			},
			expectedError: nil,
		},
		{
			name:   "error - repository error on IsDelinquent",
			loanID: 5,
			setupMocks: func(mockRepo *billingenginemocks.MockIsDelinquentRepository) {
				repoErr := errors.New("db error")
				mockRepo.On("IsDelinquent", mock.Anything, uint64(5)).Return(false, repoErr)
			},
			expectedOutput: usecases.IsDelinquentOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:   "success - customer is delinquent with no missed installments",
			loanID: 6,
			setupMocks: func(mockRepo *billingenginemocks.MockIsDelinquentRepository) {
				mockRepo.On("IsDelinquent", mock.Anything, uint64(6)).Return(true, nil)
				installments := []struct {
					WeekNumber int64
					Status     string
				}{
					{WeekNumber: 1, Status: "PAID"},
					{WeekNumber: 2, Status: "PAID"},
					{WeekNumber: 3, Status: "PENDING"},
				}
				mockRepo.On("GetInstallmentsForDelinquency", mock.Anything, uint64(6)).Return(installments, nil)
			},
			expectedOutput: usecases.IsDelinquentOutput{
				LoanID:       6,
				IsDelinquent: true,
				Message:      "Customer is delinquent - has 2 or more consecutive missed payments",
				Details: struct {
					MissedWeeks []int64 `json:"missed_weeks,omitempty"`
					TotalMissed int64   `json:"total_missed"`
				}{
					MissedWeeks: nil,
					TotalMissed: 0,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := billingenginemocks.NewMockIsDelinquentRepository(t)
			logger := zap.NewNop().Sugar()

			tt.setupMocks(mockRepo)

			interactor := NewIsDelinquentInteractor(IsDelinquentInteractorDependencies{
				IsDelinquentRepository: mockRepo,
				Logger:                 logger,
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
