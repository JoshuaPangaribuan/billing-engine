-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- PERFORMANCE INDEXES FOR BILLING ENGINE
-- =====================================================
-- These indexes are designed to optimize the most critical queries
-- based on analysis of repository methods and API endpoints
-- =====================================================

-- =====================================================
-- CUSTOMERS TABLE INDEXES
-- =====================================================

-- Primary key index already exists (id)
-- Email uniqueness constraint already provides index

-- =====================================================
-- LOANS TABLE INDEXES
-- =====================================================

-- 1. Customer ID + Status (for IsCustomerHasNonPaidLoan)
-- Critical for checking if customer has active loans
CREATE INDEX IF NOT EXISTS idx_loans_customer_id_status 
ON loans (customer_id, status);

-- 2. Customer ID + ID (for IsLoanBelongsToCustomer)
-- Critical for loan ownership validation
CREATE INDEX IF NOT EXISTS idx_loans_customer_id_id 
ON loans (customer_id, id);

-- 3. Status only (for loan status queries)
-- Useful for filtering by loan status
CREATE INDEX IF NOT EXISTS idx_loans_status 
ON loans (status);

-- 4. Start Date (for date-based queries)
-- Useful for loan creation date filtering
CREATE INDEX IF NOT EXISTS idx_loans_start_date 
ON loans (start_date);

-- =====================================================
-- INSTALLMENTS TABLE INDEXES
-- =====================================================

-- 1. Loan ID + Status (for GetOutstanding, GetPendingInstallments)
-- CRITICAL: Most frequently used combination
CREATE INDEX IF NOT EXISTS idx_installments_loan_id_status 
ON installments (loan_id, status);

-- 2. Loan ID + Week Number (for MakePayment, GetInstallments)
-- CRITICAL: Used for finding specific installments
CREATE INDEX IF NOT EXISTS idx_installments_loan_id_week_number 
ON installments (loan_id, week_number);

-- 3. Loan ID + Status + Week Number (for IsDelinquent)
-- CRITICAL: Used for delinquency detection with ordering
CREATE INDEX IF NOT EXISTS idx_installments_loan_id_status_week_number 
ON installments (loan_id, status, week_number DESC);

-- 4. Due Date + Status (for UpdateMissedInstallments)
-- CRITICAL: Used for updating overdue installments
CREATE INDEX IF NOT EXISTS idx_installments_due_date_status 
ON installments (due_date, status);

-- 5. Loan ID + Due Date (for date-based installment queries)
-- Useful for installment scheduling
CREATE INDEX IF NOT EXISTS idx_installments_loan_id_due_date 
ON installments (loan_id, due_date);

-- 6. Status only (for status-based filtering)
-- Useful for bulk status operations
CREATE INDEX IF NOT EXISTS idx_installments_status 
ON installments (status);

-- =====================================================
-- PAYMENTS TABLE INDEXES
-- =====================================================

-- 1. Installment ID (for payment lookup)
-- CRITICAL: Used for payment records
CREATE INDEX IF NOT EXISTS idx_payments_installment_id 
ON payments (installment_id);

-- 2. Paid At (for payment date queries)
-- Useful for payment history and reporting
CREATE INDEX IF NOT EXISTS idx_payments_paid_at 
ON payments (paid_at);

-- =====================================================
-- COMPOSITE INDEXES FOR COMPLEX QUERIES
-- =====================================================

-- 1. Installments: Loan ID + Status + Due Date
-- For complex outstanding calculations
CREATE INDEX IF NOT EXISTS idx_installments_loan_status_due_date 
ON installments (loan_id, status, due_date);

-- 2. Installments: Status + Due Date + Loan ID
-- For bulk status updates
CREATE INDEX IF NOT EXISTS idx_installments_status_due_date_loan 
ON installments (status, due_date, loan_id);

-- =====================================================
-- PARTIAL INDEXES FOR OPTIMIZATION
-- =====================================================

-- 1. Only PENDING installments (most queried status)
CREATE INDEX IF NOT EXISTS idx_installments_pending_only 
ON installments (loan_id, week_number, due_date) 
WHERE status = 'PENDING';

-- 2. Only MISSED installments (for delinquency detection)
CREATE INDEX IF NOT EXISTS idx_installments_missed_only 
ON installments (loan_id, week_number) 
WHERE status = 'MISSED';

-- 3. Only DISBURSED loans (most common status)
CREATE INDEX IF NOT EXISTS idx_loans_disbursed_only 
ON loans (customer_id, id) 
WHERE status = 'DISBURSED';

-- =====================================================
-- INDEXES FOR REPORTING AND ANALYTICS
-- =====================================================

-- 1. Customer + Loan + Installment relationship
-- For complex reporting queries
CREATE INDEX IF NOT EXISTS idx_installments_loan_customer_lookup 
ON installments (loan_id, status, amount_due);

-- 2. Payment tracking
-- For payment analytics
CREATE INDEX IF NOT EXISTS idx_payments_amount_paid_at 
ON payments (amount_paid, paid_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop all performance indexes in reverse order
DROP INDEX IF EXISTS idx_payments_amount_paid_at;
DROP INDEX IF EXISTS idx_installments_loan_customer_lookup;
DROP INDEX IF EXISTS idx_loans_disbursed_only;
DROP INDEX IF EXISTS idx_installments_missed_only;
DROP INDEX IF EXISTS idx_installments_pending_only;
DROP INDEX IF EXISTS idx_installments_status_due_date_loan;
DROP INDEX IF EXISTS idx_installments_loan_status_due_date;
DROP INDEX IF EXISTS idx_payments_paid_at;
DROP INDEX IF EXISTS idx_payments_installment_id;
DROP INDEX IF EXISTS idx_installments_status;
DROP INDEX IF EXISTS idx_installments_loan_id_due_date;
DROP INDEX IF EXISTS idx_installments_due_date_status;
DROP INDEX IF EXISTS idx_installments_loan_id_status_week_number;
DROP INDEX IF EXISTS idx_installments_loan_id_week_number;
DROP INDEX IF EXISTS idx_installments_loan_id_status;
DROP INDEX IF EXISTS idx_loans_start_date;
DROP INDEX IF EXISTS idx_loans_status;
DROP INDEX IF EXISTS idx_loans_customer_id_id;
DROP INDEX IF EXISTS idx_loans_customer_id_status;

-- +goose StatementEnd 