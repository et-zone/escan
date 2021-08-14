package scan

import (
	"database/sql"

	"github.com/et-zone/escan/scan/internal"
)

type EScanDef struct {
}

func NewEScanDef() *EScanDef {
	return &EScanDef{}
}

func (this *EScanDef) ScanOne(dst interface{}, rows *sql.Rows) error {
	return internal.ScanOne(dst, rows)
}

func (this *EScanDef) ScanAll(dst interface{}, rows *sql.Rows) error {
	return internal.ScanAll(dst, rows)
}
