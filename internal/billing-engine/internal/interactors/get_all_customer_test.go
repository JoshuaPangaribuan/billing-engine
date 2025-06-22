package interactors

import (
	"context"
	"errors"
	"testing"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	billingenginemocks "github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/mocks"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestGetAllCustomerInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*billingenginemocks.MockGetAllCustomerRepository)
		expectedOutput usecases.GetAllCustomerOutput
		expectedError  error
	}{
		{
			name: "success - get all customers",
			setupMocks: func(mockRepo *billingenginemocks.MockGetAllCustomerRepository) {
				customers := []entity.Customer{
					{ID: 1, Name: "John Doe", Email: "john@example.com"},
					{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
				}
				mockRepo.On("GetAllCustomer", mock.Anything).Return(customers, nil)
			},
			expectedOutput: usecases.GetAllCustomerOutput{
				Customers: []usecases.CustomerOutput{
					{ID: 1, Name: "John Doe", Email: "john@example.com"},
					{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
				},
			},
			expectedError: nil,
		},
		{
			name: "success - empty customer list",
			setupMocks: func(mockRepo *billingenginemocks.MockGetAllCustomerRepository) {
				mockRepo.On("GetAllCustomer", mock.Anything).Return([]entity.Customer{}, nil)
			},
			expectedOutput: usecases.GetAllCustomerOutput{
				Customers: []usecases.CustomerOutput{},
			},
			expectedError: nil,
		},
		{
			name: "error - repository error",
			setupMocks: func(mockRepo *billingenginemocks.MockGetAllCustomerRepository) {
				repoErr := errors.New("db error")
				mockRepo.On("GetAllCustomer", mock.Anything).Return(nil, repoErr)
			},
			expectedOutput: usecases.GetAllCustomerOutput{},
			expectedError:  &pkgerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := billingenginemocks.NewMockGetAllCustomerRepository(t)
			logger := zap.NewNop().Sugar()
			validator := validator.New()

			tt.setupMocks(mockRepo)

			interactor := NewGetAllCustomerInteractor(GetAllCustomerInteractorDependencies{
				GetAllCustomerRepository: mockRepo,
				Logger:                   logger,
				Validator:                validator,
			})

			output, err := interactor.Execute(context.Background())

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
