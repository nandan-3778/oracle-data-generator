package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Products struct {
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
	Category    string  `db:"category"`
}

var _ generator.Generator = Products{}

func (u Products) CSVHeaders() string {
	return "name,description,price,category"
}

func (u Products) CSVColumnMapping() string {
	return "(name CHAR(100), description CHAR(32000), price INTEGER, category CHAR(100))"
}

func (u Products) Table() string {
	return "products"
}

func (u Products) FakeRecord() (string, int) {
	record := fmt.Sprintf("%s,%s,%f,%s", gofakeit.StreetName(), gofakeit.Sentence(80), gofakeit.Price(10.0, 100.0), gofakeit.Sentence(2))
	return record, len(record)
}
