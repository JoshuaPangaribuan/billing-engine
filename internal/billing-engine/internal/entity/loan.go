package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type LoanStatus string

const (
	UNKNOWN_LOAN_STATUS LoanStatus = "UNKNOWN"
	LOAN_DISBURSED      LoanStatus = "DISBURSED"
	LOAN_PAID           LoanStatus = "PAID"
)

type Loan struct {
	ID              uint64          `json:"id"`
	CustomerID      uint64          `json:"customer_id"`
	PrincipalAmount decimal.Decimal `json:"principal_amount"`
	InterestRate    decimal.Decimal `json:"interest_rate"`
	TermWeeks       int64           `json:"term_weeks"`
	StartDate       time.Time       `json:"start_date"`
	Status          LoanStatus      `json:"status"`
}

// For simplicity, i use a fixed loan amount and interest rate
// and term weeks, the status is always INVESTED
func NewDisbursedLoan(customerID uint64) *Loan {
	return &Loan{
		CustomerID:      customerID,
		PrincipalAmount: decimal.NewFromUint64(5000000),
		InterestRate:    decimal.NewFromFloat(0.1),
		TermWeeks:       50,
		StartDate:       time.Now(),
		Status:          LOAN_DISBURSED,
	}
}
