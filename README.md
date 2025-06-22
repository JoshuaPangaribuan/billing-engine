# Billing Engine

A loan billing system that manages loan schedules, tracks outstanding balances, and monitors customer delinquency status.

## Overview

The Billing Engine provides:
- **Loan Schedule Generation**: Creates 50-week payment schedules for Rp 5,000,000 loans with 10% annual interest
- **Outstanding Balance Tracking**: Monitors remaining loan amounts as customers make payments
- **Delinquency Detection**: Identifies customers who miss 2 consecutive payments
- **Payment Processing**: Handles weekly repayments and catch-up payments for missed installments

### Architecture

The project follows **Clean Architecture** principles with a modular structure:

```
internal/
├── billing-engine/          # Main billing engine module
│   ├── internal/
│   │   ├── entity/          # Domain entities (Customer, Loan, Installment)
│   │   ├── usecases/        # Use case interfaces and contracts
│   │   ├── interactors/     # Business logic implementation
│   │   ├── gateway/         # External interfaces
│   │   │   ├── delivery/    # HTTP handlers and routing
│   │   │   └── repository/  # Data access layer
│   │   └── mocks/           # Generated mocks for testing
│   └── module.go            # Module initialization and dependency injection
├── app/                     # Application configuration and setup
└── pkg/                     # Shared packages and utilities
    ├── pkgerror/           # Error handling utilities
    ├── pkguid/             # ID generation (Snowflake)
    ├── pkgsql/             # Database utilities
    └── pkghttp/            # HTTP utilities
```

**Key Architectural Components**:
- **Entities**: Core business objects (Customer, Loan, Installment)
- **Use Cases**: Business logic contracts and interfaces
- **Interactors**: Implementation of business rules and workflows
- **Gateways**: External interface adapters (HTTP, Database)
- **Repository Pattern**: Data access abstraction layer
- **Dependency Injection**: Clean separation of concerns

## Key Features

### Customer Management
- **Customer Registration**: Create new customers with name and email validation
- **Customer Listing**: Retrieve all customer information

### Loan Management
- **Loan Creation**: Create new loans with automatic installment schedule generation
- **Loan Validation**: Prevent customers from having multiple unpaid loans simultaneously
- **Installment Tracking**: View detailed installment schedules with due dates and payment status

### Payment Processing
- **Weekly Payments**: Process payments for specific week numbers
- **Payment Validation**: Ensure payments match exact installment amounts and validate customer ownership
- **Customer-Loan Validation**: Verify customer exists and loan belongs to the customer before processing payments
- **Payment Status Tracking**: Monitor paid, missed, and pending installments

### Financial Tracking
- **Outstanding Balance Calculation**: Real-time calculation of remaining loan amounts
- **Payment History**: Detailed breakdown of total paid, missed, and outstanding amounts
- **Installment Status Monitoring**: Track individual installment payment status

### Delinquency Management
- **Delinquency Detection**: Automatically identify customers with 2+ consecutive missed payments
- **Delinquency Reporting**: Provide detailed reports with missed week numbers and total missed payments
- **Customer-Loan Relationship Validation**: Ensure proper ownership verification

## API Endpoints

### Customer Management
- `POST /customer` - Create a new customer
- `GET /customers` - Get all customer information

### Loan Management
- `POST /loan` - Create a new loan for a customer
- `GET /loan/:loan_id/installments` - Get installment schedule for a specific loan

### Billing Operations
- `GET /customer/:customer_id/loan/:loan_id/outstanding` - Get outstanding balance for a specific customer and loan
- `GET /loan/:loan_id/delinquent` - Check if a loan is delinquent

### Payment Operations
- `POST /loan/payment` - Process a payment for a specific loan installment
  - **Request Body**: 
    ```json
    {
      "customer_id": 1002,
      "loan_id": 2002,
      "week_number": 4,
      "amount": "110000.00"
    }
    ```
  - **Validation**: 
    - Customer must exist
    - Loan must belong to the specified customer
    - Payment amount must match the installment amount due
    - Week number must be valid for the loan

## Disclaimer

**Note**: This implementation uses hardcoded values for loan parameters:
- Principal Amount: Rp 5,000,000 (hardcoded in `NewDisbursedLoan` function)
- Interest Rate: 10% flat rate (hardcoded as 0.1)
- Term: 50 weeks (hardcoded)

These values are fixed for simplicity and cannot be configured per loan request.

**Additional Notes**:
- Negative case handling is intentionally simplified to focus on core requirements
- Error handling and edge cases are kept minimal for development simplicity
- The implementation prioritizes functionality over comprehensive error management

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.x (for development)

### Quick Start
1. **Start the application with database migration:**
   ```bash
   make up-docker
   ```

2. **Clean up and stop the application:**
   ```bash
   make clean-docker
   ```

3. **Restart the application:**
   ```bash
   make restart-docker
   ```

### Development Commands
- **Start application only:**
  ```bash
  make up
  ```

- **Stop application:**
  ```bash
  make clean
  ```

- **Run database migrations:**
  ```bash
  make migrate
  ```

- **Rollback database migrations:**
  ```bash
  make migrate-down
  ```


- **Generate mocks:**
  ```bash
  make mock
  ```

### Database Configuration
The application uses PostgreSQL with the following default configuration:
- **Host**: localhost / db (docker)
- **Port**: 5432
- **Database**: billingengine
- **Username**: root
- **Password**: rootpassword

## Testing

### API Testing with Postman

You can test the Billing Engine API using the provided Postman collection and environment:

1. **Import Postman Files**:
   - Import `Billing Engine.postman_collection.json` into Postman
   - Import `Billing Engine Env.postman_environment.json` into Postman

2. **Set Environment**:
   - Select "Billing Engine Env" environment in Postman
   - The base URL is configured as: `localhost:8081/billing-engine/api/v1`

3. **Test Data**:
   - The collection includes test data that corresponds to the database schema
   - Sample data is based on the performance indexes and seed data from `migration/20250622123000_create_performance_indexes.sql`
   - Test scenarios cover all major use cases: customer creation, loan management, payment processing, and delinquency detection

4. **Available Test Endpoints**:
   - **Customer Management**: Create and retrieve customers
   - **Loan Management**: Create loans and view installment schedules
   - **Payment Operations**: Process payments and check outstanding balances
   - **Delinquency Monitoring**: Check loan delinquency status

### Extending Functionality

For additional use cases, the API provides POST endpoints that can be extended:
- **Customer Creation**: `POST /customer` - Add new customers with validation
- **Loan Creation**: `POST /loan` - Create new loans with automatic installment generation
- **Payment Processing**: Additional payment endpoints can be added for specific business requirements

The modular architecture makes it easy to add new endpoints and business logic while maintaining clean separation of concerns.
