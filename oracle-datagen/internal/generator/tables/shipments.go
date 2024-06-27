package tables

import (
	"fmt"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"github.com/brianvoe/gofakeit"
)

type Shipments struct {
	OrderID      int    `db:"order_id"`
	TrackingNo   string `db:"tracking_no"`
	Carrier      string `db:"carrier"`
	ShipDate     string `db:"ship_date"`
	DeliveryDate string `db:"delivery_date"`
	Description  string `db:"description"`
}

var _ generator.Generator = Shipments{}

func (u Shipments) CSVHeaders() string {
	return "order_id,tracking_no,carrier,ship_date,delivery_date,description"
}

func (u Shipments) CSVColumnMapping() string {
	return "(order_id INTEGER,tracking_no CHAR(150),carrier CHAR(150),ship_date CHAR(150),delivery_date CHAR(150),description CHAR(32000))"
}

func (u Shipments) Table() string {
	return "shipments"
}

func (u Shipments) FakeRecord() (string, int) {
	record := fmt.Sprintf("%d,%s,%s,%s,%s,%s,", gofakeit.Number(1, 1000), gofakeit.StreetName(), gofakeit.StreetName(), gofakeit.Date().String(), gofakeit.Date().String(), gofakeit.Sentence(80))
	return record, len(record)
}
