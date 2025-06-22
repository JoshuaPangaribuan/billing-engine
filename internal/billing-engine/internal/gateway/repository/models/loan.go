package models

import (
	"database/sql"
	"database/sql/driver"

	"github.com/shopspring/decimal"
)

type Loan struct {
	ID              sql.NullInt64   `json:"id"`
	CustomerID      sql.NullInt64   `json:"customer_id"`
	PrincipalAmount decimal.Decimal `json:"principal_amount"`
	InterestRate    decimal.Decimal `json:"interest_rate"`
	TermWeeks       sql.NullInt64   `json:"term_weeks"`
	StartDate       sql.NullTime    `json:"start_date"`
	Status          sql.NullString  `json:"status"`
}

func (l *Loan) Columns() []any {
	return []any{
		"id",
		"customer_id",
		"principal",
		"annual_rate",
		"term_weeks",
		"start_date",
		"status",
	}
}

func (l *Loan) StringColumns() []string {
	vals := make([]string, len(l.Columns()))
	for i, col := range l.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (l *Loan) Values() []any {
	return []any{
		&l.ID,
		&l.CustomerID,
		&l.PrincipalAmount,
		&l.InterestRate,
		&l.TermWeeks,
		&l.StartDate,
		&l.Status,
	}
}

func (l Loan) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(l.Values()))
	for i, v := range l.Values() {
		vals[i] = v
	}

	return vals
}

func (l Loan) MappedValues() map[string]driver.Value {
	return map[string]driver.Value{
		"id":          l.ID.Int64,
		"customer_id": l.CustomerID.Int64,
		"principal":   l.PrincipalAmount,
		"annual_rate": l.InterestRate,
		"term_weeks":  l.TermWeeks.Int64,
		"start_date":  l.StartDate.Time,
		"status":      l.Status.String,
	}
}
