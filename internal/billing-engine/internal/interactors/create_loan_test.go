package interactors

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	billingenginemocks "github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/mocks"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgmocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateLoanInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		customerID     uint64
		setupMocks     func(*billingenginemocks.MockCreateLoanRepository, *pkgmocks.MockSnowflake)
		expectedOutput usecases.CreateLoanOutput
		expectedError  error
	}{
		{
			name:       "success - loan created successfully",
			customerID: 123,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(123)).Return(true, nil)
				mockRepo.On("IsCustomerHasNonPaidLoan", mock.Anything, uint64(123)).Return(false, nil)
				mockSnowflake.On("Generate").Return(uint64(999))

				loan := entity.NewDisbursedLoan(123)
				loan.ID = 999
				createdLoan := *loan
				mockRepo.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan entity.Loan) bool {
					return loan.CustomerID == 123 && loan.ID == 999 && loan.Status == entity.LOAN_DISBURSED
				})).Return(createdLoan, nil)
				mockRepo.On("CreateInstallmentFromLoan", mock.Anything, mock.MatchedBy(func(loan *entity.Loan) bool {
					return loan.CustomerID == 123 && loan.ID == 999 && loan.Status == entity.LOAN_DISBURSED
				})).Return(true, nil)
			},
			expectedOutput: func() usecases.CreateLoanOutput {
				loan := entity.NewDisbursedLoan(123)
				loan.ID = 999
				return usecases.CreateLoanOutput{
					ID:              999,
					CustomerID:      123,
					PrincipalAmount: loan.PrincipalAmount.String(),
					InterestRate:    loan.InterestRate.String(),
					TermWeeks:       loan.TermWeeks,
					StartDate:       loan.StartDate.Format(time.RFC3339),
					Status:          string(loan.Status),
				}
			}(),
			expectedError: nil,
		},
		{
			name:       "error - customer not found",
			customerID: 124,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(124)).Return(false, nil)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - customer has non paid loan",
			customerID: 125,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(125)).Return(true, nil)
				mockRepo.On("IsCustomerHasNonPaidLoan", mock.Anything, uint64(125)).Return(true, nil)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on IsCustomerExist",
			customerID: 126,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				repoErr := errors.New("db error")
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(126)).Return(false, repoErr)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on IsCustomerHasNonPaidLoan",
			customerID: 127,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(127)).Return(true, nil)
				repoErr := errors.New("db error")
				mockRepo.On("IsCustomerHasNonPaidLoan", mock.Anything, uint64(127)).Return(false, repoErr)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on CreateLoan",
			customerID: 128,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(128)).Return(true, nil)
				mockRepo.On("IsCustomerHasNonPaidLoan", mock.Anything, uint64(128)).Return(false, nil)
				mockSnowflake.On("Generate").Return(uint64(888))
				loan := entity.NewDisbursedLoan(128)
				loan.ID = 888
				repoErr := errors.New("db error")
				mockRepo.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan entity.Loan) bool {
					return loan.CustomerID == 128 && loan.ID == 888 && loan.Status == entity.LOAN_DISBURSED
				})).Return(entity.Loan{}, repoErr)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - repository error on CreateInstallmentFromLoan",
			customerID: 129,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(129)).Return(true, nil)
				mockRepo.On("IsCustomerHasNonPaidLoan", mock.Anything, uint64(129)).Return(false, nil)
				mockSnowflake.On("Generate").Return(uint64(777))
				loan := entity.NewDisbursedLoan(129)
				loan.ID = 777
				createdLoan := *loan
				mockRepo.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan entity.Loan) bool {
					return loan.CustomerID == 129 && loan.ID == 777 && loan.Status == entity.LOAN_DISBURSED
				})).Return(createdLoan, nil)
				repoErr := errors.New("db error")
				mockRepo.On("CreateInstallmentFromLoan", mock.Anything, mock.MatchedBy(func(loan *entity.Loan) bool {
					return loan.CustomerID == 129 && loan.ID == 777 && loan.Status == entity.LOAN_DISBURSED
				})).Return(false, repoErr)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name:       "error - failed to create installment from loan",
			customerID: 130,
			setupMocks: func(mockRepo *billingenginemocks.MockCreateLoanRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				mockRepo.On("IsCustomerExist", mock.Anything, uint64(130)).Return(true, nil)
				mockRepo.On("IsCustomerHasNonPaidLoan", mock.Anything, uint64(130)).Return(false, nil)
				mockSnowflake.On("Generate").Return(uint64(666))
				loan := entity.NewDisbursedLoan(130)
				loan.ID = 666
				createdLoan := *loan
				mockRepo.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan entity.Loan) bool {
					return loan.CustomerID == 130 && loan.ID == 666 && loan.Status == entity.LOAN_DISBURSED
				})).Return(createdLoan, nil)
				mockRepo.On("CreateInstallmentFromLoan", mock.Anything, mock.MatchedBy(func(loan *entity.Loan) bool {
					return loan.CustomerID == 130 && loan.ID == 666 && loan.Status == entity.LOAN_DISBURSED
				})).Return(false, nil)
			},
			expectedOutput: usecases.CreateLoanOutput{},
			expectedError:  &pkgerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := billingenginemocks.NewMockCreateLoanRepository(t)
			mockSnowflake := pkgmocks.NewMockSnowflake(t)
			logger := zap.NewNop().Sugar()

			tt.setupMocks(mockRepo, mockSnowflake)

			interactor := NewCreateLoanInteractor(CreateLoanInteractorDependencies{
				CreateLoanRepository: mockRepo,
				Logger:               logger,
				SnowflakeGen:         mockSnowflake,
			})

			output, err := interactor.Execute(context.Background(), tt.customerID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			mockRepo.AssertExpectations(t)
			mockSnowflake.AssertExpectations(t)
		})
	}
}
