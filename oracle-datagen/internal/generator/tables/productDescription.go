package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type ProductDescription struct {
	Name        string `db:"name"`
	Description string `db:"description"`
}

var _ generator.Generator = ProductDescription{}

func (u ProductDescription) CSVHeaders() string {
	return "name,description"
}

func (u ProductDescription) CSVColumnMapping() string {
	return "(name CHAR(100), description CHAR(32000))"
}

func (u ProductDescription) Table() string {
	return "ProductDescription"
}

func (u ProductDescription) FakeRecord() (string, int) {
	record := fmt.Sprintf("%s,%s", gofakeit.Name(), gofakeit.Sentence(100))
	return record, len(record)
}
