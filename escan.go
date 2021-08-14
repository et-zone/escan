package escan

import (
	"database/sql"

	"github.com/et-zone/escan/scan"
)

//explanation
//use db rows ,all of the db field can not be null,so you must set value every field
//

type EScan interface {
	ScanOne(dst interface{}, rows *sql.Rows) error
	ScanAll(dst interface{}, rows *sql.Rows) error
}

func NewEscan() EScan {
	return scan.NewEScanDef()
}
