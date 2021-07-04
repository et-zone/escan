package escan

import (
	"database/sql"

	"github.con/et-zone/escan/scan"
)

type EScan interface {
	ScanOne(dst interface{}, rows *sql.Rows) error
	ScanAll(dst interface{}, rows *sql.Rows) error
}

func NewEscan() EScan {
	return scan.NewEScanDef()
}
