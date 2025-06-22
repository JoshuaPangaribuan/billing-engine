-- +goose Up
CREATE TABLE IF NOT EXISTS customers (
  id BIGINT NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS loans (
  id BIGINT NOT NULL PRIMARY KEY,
  customer_id BIGINT NOT NULL, -- FK to customer.id
  principal DECIMAL(10, 2) NOT NULL, 
  annual_rate DECIMAL(5, 4) NOT NULL,
  term_weeks INT NOT NULL,
  start_date DATE NOT NULL, -- RFC 3339
  status VARCHAR(20) NOT NULL CHECK (status IN ('DISBURSED', 'PAID'))
);

CREATE TABLE IF NOT EXISTS installments (
    id BIGSERIAL PRIMARY KEY,  
    loan_id BIGINT NOT NULL,                        
    week_number INT NOT NULL,                                
    due_date DATE NOT NULL,                                  
    amount_due DECIMAL(18,2) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'PAID', 'MISSED')),
    UNIQUE (loan_id, week_number) -- Satu installment per minggu per loan
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    installment_id BIGINT NOT NULL,                 
    paid_at TIMESTAMP NOT NULL,                               
    amount_paid DECIMAL(18,2) NOT NULL                       
);

-- +goose Down
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS installments;
DROP TABLE IF EXISTS payments;