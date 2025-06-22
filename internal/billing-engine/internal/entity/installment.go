package entity

type InstallmentStatus string

const (
	INSTALLMENT_UNKNOWN InstallmentStatus = "UNKNOWN"
	INSTALLMENT_PENDING InstallmentStatus = "PENDING"
	INSTALLMENT_PAID    InstallmentStatus = "PAID"
	INSTALLMENT_MISSED  InstallmentStatus = "MISSED"
)

type Installment struct {
	ID         uint64            `json:"id"`
	LoanID     uint64            `json:"loan_id"`
	WeekNumber int64             `json:"week_number"`
	DueDate    string            `json:"due_date"`
	AmountDue  string            `json:"amount_due"`
	Status     InstallmentStatus `json:"status"`
}
