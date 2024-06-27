package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Payment struct {
	OrderID       int    `db:"order_id"`
	Amount        int    `db:"amount"`
	PaymentMethod string `db:"payment_method"`
	PaymentDate   string `db:"payment_date"`
	Description   string `db:"description"`
}

var _ generator.Generator = Payment{}

func (p Payment) CSVHeaders() string {
	return "order_id, amount, payment_method, payment_date,description"
}

func (p Payment) CSVColumnMapping() string {
	return "(order_id INTEGER, amount INTERGER, payment_method CHAR(150), payment_date CHAR(150),description CHAR(32000)"
}

func (p Payment) Table() string {
	return "payments"
}

func (p Payment) FakeRecord() (string, int) {
	record := fmt.Sprintf("%d,%d,%s,%s,%s", gofakeit.Number(1, 100), gofakeit.Number(1, 50), gofakeit.CreditCardType(), gofakeit.Date().String(), gofakeit.Sentence(80))
	return record, len(record)
}
