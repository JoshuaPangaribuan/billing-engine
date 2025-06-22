package interactors

import (
	"context"
	"errors"
	"testing"

	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/entity"
	billingenginemocks "github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/mocks"
	"github.com/JoshuaPangaribuan/billing-engine/internal/billing-engine/internal/usecases"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgerror"
	"github.com/JoshuaPangaribuan/billing-engine/internal/pkg/pkgmocks"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateCustomerInteractor_Execute(t *testing.T) {
	tests := []struct {
		name           string
		input          usecases.CreateCustomerInput
		setupMocks     func(*billingenginemocks.MockCreateCustomerRepository, *pkgmocks.MockSnowflake)
		expectedOutput usecases.CreateCustomerOutput
		expectedError  error
	}{
		{
			name: "success - customer created successfully",
			input: usecases.CreateCustomerInput{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				expectedID := uint64(123456789)
				expectedCustomer := entity.Customer{
					ID:    expectedID,
					Name:  "John Doe",
					Email: "john.doe@example.com",
				}

				mockSnowflake.On("Generate").Return(expectedID)
				mockRepo.On("CreateCustomer", mock.Anything, mock.MatchedBy(func(customer entity.Customer) bool {
					return customer.Name == "John Doe" && customer.Email == "john.doe@example.com" && customer.ID == expectedID
				})).Return(expectedCustomer, nil)
			},
			expectedOutput: usecases.CreateCustomerOutput{
				ID:    123456789,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expectedError: nil,
		},
		{
			name: "validation error - empty name",
			input: usecases.CreateCustomerInput{
				Name:  "",
				Email: "john.doe@example.com",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.CreateCustomerOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "validation error - empty email",
			input: usecases.CreateCustomerInput{
				Name:  "John Doe",
				Email: "",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.CreateCustomerOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "validation error - invalid email format",
			input: usecases.CreateCustomerInput{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				// No mocks needed for validation error
			},
			expectedOutput: usecases.CreateCustomerOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "repository error - database error",
			input: usecases.CreateCustomerInput{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				expectedID := uint64(123456789)
				dbError := errors.New("database connection failed")

				mockSnowflake.On("Generate").Return(expectedID)
				mockRepo.On("CreateCustomer", mock.Anything, mock.MatchedBy(func(customer entity.Customer) bool {
					return customer.Name == "John Doe" && customer.Email == "john.doe@example.com" && customer.ID == expectedID
				})).Return(entity.Customer{}, dbError)
			},
			expectedOutput: usecases.CreateCustomerOutput{},
			expectedError:  &pkgerror.Error{},
		},
		{
			name: "repository error - duplicate email",
			input: usecases.CreateCustomerInput{
				Name:  "John Doe",
				Email: "existing@example.com",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				expectedID := uint64(123456789)
				duplicateError := errors.New("duplicate email")

				mockSnowflake.On("Generate").Return(expectedID)
				mockRepo.On("CreateCustomer", mock.Anything, mock.MatchedBy(func(customer entity.Customer) bool {
					return customer.Name == "John Doe" && customer.Email == "existing@example.com" && customer.ID == expectedID
				})).Return(entity.Customer{}, duplicateError)
			},
			expectedOutput: usecases.CreateCustomerOutput{},
			expectedError:  &pkgerror.Error{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := billingenginemocks.NewMockCreateCustomerRepository(t)
			mockSnowflake := pkgmocks.NewMockSnowflake(t)
			logger := zap.NewNop().Sugar()
			validator := validator.New()

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockSnowflake)

			// Create interactor
			interactor := NewCreateCustomerInteractor(CreateCustomerInteractorDependencies{
				CustomerRepository: mockRepo,
				Logger:             logger,
				Validator:          validator,
				SnowflakeGen:       mockSnowflake,
			})

			// Execute
			output, err := interactor.Execute(context.Background(), tt.input)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			// Verify all mocks were called as expected
			mockRepo.AssertExpectations(t)
			mockSnowflake.AssertExpectations(t)
		})
	}
}

func TestNewCreateCustomerInteractor(t *testing.T) {
	tests := []struct {
		name        string
		deps        CreateCustomerInteractorDependencies
		expectPanic bool
	}{
		{
			name: "success - valid dependencies",
			deps: CreateCustomerInteractorDependencies{
				CustomerRepository: billingenginemocks.NewMockCreateCustomerRepository(t),
				Logger:             zap.NewNop().Sugar(),
				Validator:          validator.New(),
				SnowflakeGen:       pkgmocks.NewMockSnowflake(t),
			},
			expectPanic: false,
		},
		{
			name: "success - nil repository (no validation on dependencies struct)",
			deps: CreateCustomerInteractorDependencies{
				CustomerRepository: nil,
				Logger:             zap.NewNop().Sugar(),
				Validator:          validator.New(),
				SnowflakeGen:       pkgmocks.NewMockSnowflake(t),
			},
			expectPanic: false,
		},
		{
			name: "success - nil logger (no validation on dependencies struct)",
			deps: CreateCustomerInteractorDependencies{
				CustomerRepository: billingenginemocks.NewMockCreateCustomerRepository(t),
				Logger:             nil,
				Validator:          validator.New(),
				SnowflakeGen:       pkgmocks.NewMockSnowflake(t),
			},
			expectPanic: false,
		},
		{
			name: "success - nil validator (no validation on dependencies struct)",
			deps: CreateCustomerInteractorDependencies{
				CustomerRepository: billingenginemocks.NewMockCreateCustomerRepository(t),
				Logger:             zap.NewNop().Sugar(),
				Validator:          nil,
				SnowflakeGen:       pkgmocks.NewMockSnowflake(t),
			},
			expectPanic: false,
		},
		{
			name: "success - nil snowflake generator (no validation on dependencies struct)",
			deps: CreateCustomerInteractorDependencies{
				CustomerRepository: billingenginemocks.NewMockCreateCustomerRepository(t),
				Logger:             zap.NewNop().Sugar(),
				Validator:          validator.New(),
				SnowflakeGen:       nil,
			},
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					NewCreateCustomerInteractor(tt.deps)
				})
			} else {
				assert.NotPanics(t, func() {
					interactor := NewCreateCustomerInteractor(tt.deps)
					assert.NotNil(t, interactor)
					assert.Equal(t, tt.deps.CustomerRepository, interactor.repository)
					assert.Equal(t, tt.deps.Logger, interactor.logger)
					assert.Equal(t, tt.deps.Validator, interactor.validator)
					assert.Equal(t, tt.deps.SnowflakeGen, interactor.snowflakeGen)
				})
			}
		})
	}
}

func TestCreateCustomerInteractor_Execute_EdgeCases(t *testing.T) {
	tests := []struct {
		name           string
		input          usecases.CreateCustomerInput
		setupMocks     func(*billingenginemocks.MockCreateCustomerRepository, *pkgmocks.MockSnowflake)
		expectedOutput usecases.CreateCustomerOutput
		expectedError  error
	}{
		{
			name: "success - with special characters in name",
			input: usecases.CreateCustomerInput{
				Name:  "José María O'Connor-Smith",
				Email: "jose.maria@example.com",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				expectedID := uint64(987654321)
				expectedCustomer := entity.Customer{
					ID:    expectedID,
					Name:  "José María O'Connor-Smith",
					Email: "jose.maria@example.com",
				}

				mockSnowflake.On("Generate").Return(expectedID)
				mockRepo.On("CreateCustomer", mock.Anything, mock.MatchedBy(func(customer entity.Customer) bool {
					return customer.Name == "José María O'Connor-Smith" && customer.Email == "jose.maria@example.com" && customer.ID == expectedID
				})).Return(expectedCustomer, nil)
			},
			expectedOutput: usecases.CreateCustomerOutput{
				ID:    987654321,
				Name:  "José María O'Connor-Smith",
				Email: "jose.maria@example.com",
			},
			expectedError: nil,
		},
		{
			name: "success - with long name and email",
			input: usecases.CreateCustomerInput{
				Name:  "Very Long Customer Name That Exceeds Normal Length Expectations",
				Email: "very.long.email.address.for.customer@very-long-domain-name.example.com",
			},
			setupMocks: func(mockRepo *billingenginemocks.MockCreateCustomerRepository, mockSnowflake *pkgmocks.MockSnowflake) {
				expectedID := uint64(111222333)
				expectedCustomer := entity.Customer{
					ID:    expectedID,
					Name:  "Very Long Customer Name That Exceeds Normal Length Expectations",
					Email: "very.long.email.address.for.customer@very-long-domain-name.example.com",
				}

				mockSnowflake.On("Generate").Return(expectedID)
				mockRepo.On("CreateCustomer", mock.Anything, mock.MatchedBy(func(customer entity.Customer) bool {
					return customer.Name == "Very Long Customer Name That Exceeds Normal Length Expectations" &&
						customer.Email == "very.long.email.address.for.customer@very-long-domain-name.example.com" &&
						customer.ID == expectedID
				})).Return(expectedCustomer, nil)
			},
			expectedOutput: usecases.CreateCustomerOutput{
				ID:    111222333,
				Name:  "Very Long Customer Name That Exceeds Normal Length Expectations",
				Email: "very.long.email.address.for.customer@very-long-domain-name.example.com",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := billingenginemocks.NewMockCreateCustomerRepository(t)
			mockSnowflake := pkgmocks.NewMockSnowflake(t)
			logger := zap.NewNop().Sugar()
			validator := validator.New()

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockSnowflake)

			// Create interactor
			interactor := NewCreateCustomerInteractor(CreateCustomerInteractorDependencies{
				CustomerRepository: mockRepo,
				Logger:             logger,
				Validator:          validator,
				SnowflakeGen:       mockSnowflake,
			})

			// Execute
			output, err := interactor.Execute(context.Background(), tt.input)

			// Assertions
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, output)
			}

			// Verify all mocks were called as expected
			mockRepo.AssertExpectations(t)
			mockSnowflake.AssertExpectations(t)
		})
	}
}
