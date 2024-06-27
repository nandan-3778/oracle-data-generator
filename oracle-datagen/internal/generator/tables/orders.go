package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Orders struct {
	UserID      int    `db:"user_id"`
	ProductID   int    `db:"product_id"`
	Quantity    int    `db:"quantity"`
	Total       int    `db:"total"`
	OrderDate   string `db:"order_date"`
	Description string `db:"description"`
}

var _ generator.Generator = Orders{}

func (u Orders) CSVHeaders() string {
	return "user_id,product_id,quantity,total,order_date,description"
}

func (u Orders) CSVColumnMapping() string {
	return "(user_id INTEGER, product_id INTEGER, quantity INTEGER, total INTEGER,order_date CHAR(100), description CHAR(32000))"
}

func (u Orders) Table() string {
	return "orders"
}

func (u Orders) FakeRecord() (string, int) {
	record := fmt.Sprintf("%d,%d,%d,%d,%s,%s", gofakeit.Number(1, 100), gofakeit.Number(1, 50), gofakeit.Number(1, 5), gofakeit.Number(1, 5), gofakeit.Date().Format("2006-01-02"), gofakeit.Sentence(80))
	return record, len(record)
}
