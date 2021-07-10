package escan

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.con/et-zone/escan/scan"
)

type EScan interface {
	ScanOne(dst interface{}, rows *sql.Rows) error
	ScanAll(dst interface{}, rows *sql.Rows) error
}

func NewEscan() EScan {
	return scan.NewEScanDef()
}

//dbKey to jsonKey
func ChToJsonByTagDB(src map[string]string, desStruct interface{}) (map[string]string, error) {
	newTagdata := map[string]string{}
	if src == nil {
		return newTagdata, errors.New("tagData is nil,the tagData neet init")
	}

	tVal := reflect.TypeOf(desStruct)

	if tVal.Kind() != reflect.Struct {
		return newTagdata, errors.New("desStruct must struct type ")
	}

	l := tVal.NumField()
	for i := 0; i < l; i++ {
		jsonName := tVal.Field(i).Tag.Get("json")
		tagName := tVal.Field(i).Tag.Get("db")
		jsonName = strings.TrimSpace(strings.Split(jsonName, ",")[0])
		if jsonName == "-" || tagName == "-" {
			continue
		}
		if jsonName == "" {
			jsonName = tVal.Field(i).Name
		}
		if tagName == "" {
			tagName = tVal.Field(i).Name
		}
		if jsonName == tagName {
			newTagdata[jsonName] = src[jsonName]
		}
		if src[tagName] != "" {
			newTagdata[jsonName] = src[tagName]
		}

	}

	return newTagdata, nil
}
