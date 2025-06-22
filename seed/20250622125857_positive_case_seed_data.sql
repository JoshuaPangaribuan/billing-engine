-- +goose Up
-- +goose StatementBegin

-- Customer 1: John Doe - Has active loan with some payments made
INSERT INTO customers (id, name, email) VALUES 
(1001, 'John Doe', 'john.doe@example.com');

-- Customer 2: Jane Smith - Has active loan with different payment pattern
INSERT INTO customers (id, name, email) VALUES 
(1002, 'Jane Smith', 'jane.smith@example.com');

-- Loan for Customer 1: Started 10 weeks ago, some payments made
INSERT INTO loans (id, customer_id, principal, annual_rate, term_weeks, start_date, status) VALUES 
(2001, 1001, 5000000.00, 0.1000, 50, '2024-01-01', 'DISBURSED');

-- Loan for Customer 2: Started 5 weeks ago, different payment pattern
INSERT INTO loans (id, customer_id, principal, annual_rate, term_weeks, start_date, status) VALUES 
(2002, 1002, 5000000.00, 0.1000, 50, '2024-02-01', 'DISBURSED');

-- Installments for Customer 1's loan (50 weeks total)
-- Week 1-5: PAID (fully paid)
-- Week 6-8: PENDING (not yet due)
-- Week 9-10: MISSED (overdue)
-- Week 11-50: PENDING (future installments)

-- Week 1-5: PAID installments
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2001, 1, '2024-01-08', 110000.00, 'PAID'),
(2001, 2, '2024-01-15', 110000.00, 'PAID'),
(2001, 3, '2024-01-22', 110000.00, 'PAID'),
(2001, 4, '2024-01-29', 110000.00, 'PAID'),
(2001, 5, '2024-02-05', 110000.00, 'PAID');

-- Week 6-8: PENDING installments (current week and upcoming)
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2001, 6, '2024-02-12', 110000.00, 'PENDING'),
(2001, 7, '2024-02-19', 110000.00, 'PENDING'),
(2001, 8, '2024-02-26', 110000.00, 'PENDING');

-- Week 9-10: MISSED installments (overdue)
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2001, 9, '2024-02-05', 110000.00, 'MISSED'),
(2001, 10, '2024-02-12', 110000.00, 'MISSED');

-- Week 11-50: Future PENDING installments
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2001, 11, '2024-03-05', 110000.00, 'PENDING'),
(2001, 12, '2024-03-12', 110000.00, 'PENDING'),
(2001, 13, '2024-03-19', 110000.00, 'PENDING'),
(2001, 14, '2024-03-26', 110000.00, 'PENDING'),
(2001, 15, '2024-04-02', 110000.00, 'PENDING'),
(2001, 16, '2024-04-09', 110000.00, 'PENDING'),
(2001, 17, '2024-04-16', 110000.00, 'PENDING'),
(2001, 18, '2024-04-23', 110000.00, 'PENDING'),
(2001, 19, '2024-04-30', 110000.00, 'PENDING'),
(2001, 20, '2024-05-07', 110000.00, 'PENDING'),
(2001, 21, '2024-05-14', 110000.00, 'PENDING'),
(2001, 22, '2024-05-21', 110000.00, 'PENDING'),
(2001, 23, '2024-05-28', 110000.00, 'PENDING'),
(2001, 24, '2024-06-04', 110000.00, 'PENDING'),
(2001, 25, '2024-06-11', 110000.00, 'PENDING'),
(2001, 26, '2024-06-18', 110000.00, 'PENDING'),
(2001, 27, '2024-06-25', 110000.00, 'PENDING'),
(2001, 28, '2024-07-02', 110000.00, 'PENDING'),
(2001, 29, '2024-07-09', 110000.00, 'PENDING'),
(2001, 30, '2024-07-16', 110000.00, 'PENDING'),
(2001, 31, '2024-07-23', 110000.00, 'PENDING'),
(2001, 32, '2024-07-30', 110000.00, 'PENDING'),
(2001, 33, '2024-08-06', 110000.00, 'PENDING'),
(2001, 34, '2024-08-13', 110000.00, 'PENDING'),
(2001, 35, '2024-08-20', 110000.00, 'PENDING'),
(2001, 36, '2024-08-27', 110000.00, 'PENDING'),
(2001, 37, '2024-09-03', 110000.00, 'PENDING'),
(2001, 38, '2024-09-10', 110000.00, 'PENDING'),
(2001, 39, '2024-09-17', 110000.00, 'PENDING'),
(2001, 40, '2024-09-24', 110000.00, 'PENDING'),
(2001, 41, '2024-10-01', 110000.00, 'PENDING'),
(2001, 42, '2024-10-08', 110000.00, 'PENDING'),
(2001, 43, '2024-10-15', 110000.00, 'PENDING'),
(2001, 44, '2024-10-22', 110000.00, 'PENDING'),
(2001, 45, '2024-10-29', 110000.00, 'PENDING'),
(2001, 46, '2024-11-05', 110000.00, 'PENDING'),
(2001, 47, '2024-11-12', 110000.00, 'PENDING'),
(2001, 48, '2024-11-19', 110000.00, 'PENDING'),
(2001, 49, '2024-11-26', 110000.00, 'PENDING'),
(2001, 50, '2024-12-03', 110000.00, 'PENDING');

-- Installments for Customer 2's loan (50 weeks total)
-- Week 1-3: PAID (fully paid)
-- Week 4-5: PENDING (not yet due)
-- Week 6-50: PENDING (future installments)

-- Week 1-3: PAID installments
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2002, 1, '2024-02-08', 110000.00, 'PAID'),
(2002, 2, '2024-02-15', 110000.00, 'PAID'),
(2002, 3, '2024-02-22', 110000.00, 'PAID');

-- Week 4-5: PENDING installments (current week and upcoming)
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2002, 4, '2024-02-29', 110000.00, 'PENDING'),
(2002, 5, '2024-03-07', 110000.00, 'PENDING');

-- Week 6-50: Future PENDING installments
INSERT INTO installments (loan_id, week_number, due_date, amount_due, status) VALUES
(2002, 6, '2024-03-14', 110000.00, 'PENDING'),
(2002, 7, '2024-03-21', 110000.00, 'PENDING'),
(2002, 8, '2024-03-28', 110000.00, 'PENDING'),
(2002, 9, '2024-04-04', 110000.00, 'PENDING'),
(2002, 10, '2024-04-11', 110000.00, 'PENDING'),
(2002, 11, '2024-04-18', 110000.00, 'PENDING'),
(2002, 12, '2024-04-25', 110000.00, 'PENDING'),
(2002, 13, '2024-05-02', 110000.00, 'PENDING'),
(2002, 14, '2024-05-09', 110000.00, 'PENDING'),
(2002, 15, '2024-05-16', 110000.00, 'PENDING'),
(2002, 16, '2024-05-23', 110000.00, 'PENDING'),
(2002, 17, '2024-05-30', 110000.00, 'PENDING'),
(2002, 18, '2024-06-06', 110000.00, 'PENDING'),
(2002, 19, '2024-06-13', 110000.00, 'PENDING'),
(2002, 20, '2024-06-20', 110000.00, 'PENDING'),
(2002, 21, '2024-06-27', 110000.00, 'PENDING'),
(2002, 22, '2024-07-04', 110000.00, 'PENDING'),
(2002, 23, '2024-07-11', 110000.00, 'PENDING'),
(2002, 24, '2024-07-18', 110000.00, 'PENDING'),
(2002, 25, '2024-07-25', 110000.00, 'PENDING'),
(2002, 26, '2024-08-01', 110000.00, 'PENDING'),
(2002, 27, '2024-08-08', 110000.00, 'PENDING'),
(2002, 28, '2024-08-15', 110000.00, 'PENDING'),
(2002, 29, '2024-08-22', 110000.00, 'PENDING'),
(2002, 30, '2024-08-29', 110000.00, 'PENDING'),
(2002, 31, '2024-09-05', 110000.00, 'PENDING'),
(2002, 32, '2024-09-12', 110000.00, 'PENDING'),
(2002, 33, '2024-09-19', 110000.00, 'PENDING'),
(2002, 34, '2024-09-26', 110000.00, 'PENDING'),
(2002, 35, '2024-10-03', 110000.00, 'PENDING'),
(2002, 36, '2024-10-10', 110000.00, 'PENDING'),
(2002, 37, '2024-10-17', 110000.00, 'PENDING'),
(2002, 38, '2024-10-24', 110000.00, 'PENDING'),
(2002, 39, '2024-10-31', 110000.00, 'PENDING'),
(2002, 40, '2024-11-07', 110000.00, 'PENDING'),
(2002, 41, '2024-11-14', 110000.00, 'PENDING'),
(2002, 42, '2024-11-21', 110000.00, 'PENDING'),
(2002, 43, '2024-11-28', 110000.00, 'PENDING'),
(2002, 44, '2024-12-05', 110000.00, 'PENDING'),
(2002, 45, '2024-12-12', 110000.00, 'PENDING'),
(2002, 46, '2024-12-19', 110000.00, 'PENDING'),
(2002, 47, '2024-12-26', 110000.00, 'PENDING'),
(2002, 48, '2025-01-02', 110000.00, 'PENDING'),
(2002, 49, '2025-01-09', 110000.00, 'PENDING'),
(2002, 50, '2025-01-16', 110000.00, 'PENDING');

-- Payments for Customer 1's loan (for the PAID installments)
INSERT INTO payments (installment_id, paid_at, amount_paid) VALUES
(1, '2024-01-08 10:30:00', 110000.00),  -- Week 1 payment
(2, '2024-01-15 14:20:00', 110000.00),  -- Week 2 payment
(3, '2024-01-22 09:15:00', 110000.00),  -- Week 3 payment
(4, '2024-01-29 16:45:00', 110000.00),  -- Week 4 payment
(5, '2024-02-05 11:30:00', 110000.00);  -- Week 5 payment

-- Payments for Customer 2's loan (for the PAID installments)
INSERT INTO payments (installment_id, paid_at, amount_paid) VALUES
(56, '2024-02-08 13:20:00', 110000.00), -- Week 1 payment
(57, '2024-02-15 10:45:00', 110000.00), -- Week 2 payment
(58, '2024-02-22 15:30:00', 110000.00); -- Week 3 payment

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete all test data in reverse order
DELETE FROM payments WHERE installment_id IN (
    SELECT id FROM installments WHERE loan_id IN (2001, 2002)
);

DELETE FROM installments WHERE loan_id IN (2001, 2002);

DELETE FROM loans WHERE id IN (2001, 2002);

DELETE FROM customers WHERE id IN (1001, 1002);

-- +goose StatementEnd
