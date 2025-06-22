package models

import (
	"database/sql"
	"database/sql/driver"
)

type Customer struct {
	ID    sql.NullInt64  `json:"id"`
	Name  sql.NullString `json:"name"`
	Email sql.NullString `json:"email"`
}

func (c *Customer) Columns() []any {
	return []any{
		"id",
		"name",
		"email",
	}
}

func (c *Customer) StringColumns() []string {
	vals := make([]string, len(c.Columns()))
	for i, col := range c.Columns() {
		c, ok := col.(string)
		if ok {
			vals[i] = c
		}
	}

	return vals
}

func (c *Customer) Values() []any {
	return []any{
		&c.ID,
		&c.Name,
		&c.Email,
	}
}

func (c Customer) DriverValues() []driver.Value {
	vals := make([]driver.Value, len(c.Values()))
	for i, v := range c.Values() {
		vals[i] = v
	}

	return vals
}

func (c Customer) MappedValues() map[string]driver.Value {
	return map[string]driver.Value{
		"id":    c.ID,
		"name":  c.Name.String,
		"email": c.Email.String,
	}
}
