package models

import (
	"database/sql"
	"database/sql/driver"
)

type Payment struct {
	ID            sql.NullInt64  `json:"id"`
	InstallmentID sql.NullInt64  `json:"installment_id"`
	PaidAt        sql.NullTime   `json:"paid_at"`
	AmountPaid    sql.NullString `json:"amount_paid"`
}

func (p *Payment) Columns() []any {
	return []any{
		"id",
		"installment_id",
		"paid_at",
		"amount_paid",
	}
}

func (p *Payment) StringColumns() []string {
	vals := make([]string, len(p.Columns()))
	for i, col := range p.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (p *Payment) Values() []any {
	return []any{
		&p.ID,
		&p.InstallmentID,
		&p.PaidAt,
		&p.AmountPaid,
	}
}

func (p Payment) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(p.Values()))
	for i, v := range p.Values() {
		vals[i] = v
	}

	return vals
}

func (p Payment) MappedValues() map[string]driver.Value {
	return map[string]driver.Value{
		"id":             p.ID.Int64,
		"installment_id": p.InstallmentID.Int64,
		"paid_at":        p.PaidAt.Time,
		"amount_paid":    p.AmountPaid.String,
	}
}
