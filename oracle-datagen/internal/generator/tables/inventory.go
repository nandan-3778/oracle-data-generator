package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Inventory struct {
	ProductID   int    `db:"product_id"`
	Quantity    int    `db:"quantity"`
	ReorderLvl  int    `db:"reorder_lvl"`
	Description string `db:"description"`
}

var _ generator.Generator = Inventory{}

func (u Inventory) CSVHeaders() string {
	return "product_id,quantity,reorder_lvl,description"
}

func (u Inventory) CSVColumnMapping() string {
	return "(product_id INTEGER,quantity INTEGER,reorder_lvl INTEGER,description CHAR(32000))"
}

func (u Inventory) Table() string {
	return "Inventory"
}

func (u Inventory) FakeRecord() (string, int) {
	record := fmt.Sprintf("%d,%d.%d,%s", gofakeit.Number(1, 50), gofakeit.Number(1, 100), gofakeit.Number(1, 50), gofakeit.Sentence(80))
	return record, len(record)
}
