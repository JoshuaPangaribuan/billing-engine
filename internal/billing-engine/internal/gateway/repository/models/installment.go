package models

import (
	"database/sql"
	"database/sql/driver"
)

type Installment struct {
	ID         sql.NullInt64  `json:"id"`
	LoanID     sql.NullInt64  `json:"loan_id"`
	WeekNumber sql.NullInt64  `json:"week_number"`
	DueDate    sql.NullString `json:"due_date"`
	AmountDue  sql.NullString `json:"amount_due"`
	Status     sql.NullString `json:"status"`
}

func (i *Installment) Columns() []any {
	return []any{
		"id",
		"loan_id",
		"week_number",
		"due_date",
		"amount_due",
		"status",
	}
}

func (i *Installment) StringColumns() []string {
	vals := make([]string, len(i.Columns()))
	for j, col := range i.Columns() {
		c, ok := col.(string)
		if ok {
			vals[j] = c
		}
	}

	return vals
}

func (i *Installment) Values() []any {
	return []any{
		&i.ID,
		&i.LoanID,
		&i.WeekNumber,
		&i.DueDate,
		&i.AmountDue,
		&i.Status,
	}
}

func (i Installment) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(i.Values()))
	for j, v := range i.Values() {
		vals[j] = v
	}

	return vals
}

func (i Installment) MappedValues() map[string]driver.Value {
	return map[string]driver.Value{
		"id":          i.ID.Int64,
		"loan_id":     i.LoanID.Int64,
		"week_number": i.WeekNumber.Int64,
		"due_date":    i.DueDate.String,
		"amount_due":  i.AmountDue.String,
		"status":      i.Status.String,
	}
}
