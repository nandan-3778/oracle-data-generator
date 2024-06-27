package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Suppliers struct {
	Name        string `db:"name"`
	Address     string `db:"address"`
	City        string `db:"city"`
	State       string `db:"state"`
	ZipCode     string `db:"zip_code"`
	Country     string `db:"country"`
	Phone       string `db:"phone"`
	Description string `db:"description"`
}

var _ generator.Generator = Suppliers{}

func (u Suppliers) CSVHeaders() string {
	return "name,address,city,state,zip_code,country,phone,description"
}

func (u Suppliers) CSVColumnMapping() string {
	return "(name CHAR(300),address CHAR(300),city CHAR(100),state CHAR(100),zip_code CHAR(100),country CHAR(100),phone CHAR(100),description CHAR(32000))"
}

func (u Suppliers) Table() string {
	return "suppliers"
}

// func (u Suppliers) FakeRecord() string {
// 	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,", gofakeit.Company(), gofakeit.Address().Address, gofakeit.City(), gofakeit.State(), gofakeit.Zip(), gofakeit.Country(), gofakeit.Phone(), gofakeit.Sentence(80))
// }

func (u Suppliers) FakeRecord() (string, int) {
	record := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,", gofakeit.Company(), gofakeit.Address().Address, gofakeit.City(), gofakeit.State(), gofakeit.Zip(), gofakeit.Country(), gofakeit.Phone(), gofakeit.Sentence(80))
	return record, len(record)
}
