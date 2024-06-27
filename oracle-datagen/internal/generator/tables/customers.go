package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Customer struct {
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Phone       string `db:"phone"`
	Description string `db:"description"`
}

var _ generator.Generator = Customer{}

func (u Customer) CSVHeaders() string {
	return "first_name,last_name,email,phone,description"
}

func (u Customer) CSVColumnMapping() string {
	return "(first_name CHAR(100), last_name CHAR(100), email CHAR(100), phone INTEGER, description CHAR(32000))"
}

func (u Customer) Table() string {
	return "customers"
}

func (u Customer) FakeRecord() (string, int) {
	record := fmt.Sprintf("%s,%s,%s,%s,%s", gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), gofakeit.Phone(), gofakeit.Sentence(70))
	return record, len(record)
}
