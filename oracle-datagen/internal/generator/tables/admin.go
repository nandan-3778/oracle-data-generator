package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Admin struct {
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string `db:"email"`
	Age         int    `db:"age"`
	Description string `db:"description"`
}

var _ generator.Generator = Admin{}

func (u Admin) CSVHeaders() string {
	return "first_name,last_name,email,age,description"
}

func (u Admin) CSVColumnMapping() string {
	return "(first_name CHAR(100), last_name CHAR(100), email CHAR(100), age INTEGER, description CHAR(32000))"
}

func (u Admin) Table() string {
	return "admin"
}

func (u Admin) FakeRecord() (string, int) {
	record := fmt.Sprintf("%s,%s,%s,%d,%s", gofakeit.FirstName(), gofakeit.LastName(), gofakeit.Email(), gofakeit.Number(18, 75), gofakeit.Sentence(70))
	return record, len(record)
}
