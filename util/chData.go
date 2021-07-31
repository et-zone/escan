package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	DEFAULT_NUM = 1
	DEFAULT_STR = "1"
)

func getTagName(f reflect.StructField, tag string) string {
	val, ok := f.Tag.Lookup(tag)
	fieldTag := ""
	if !ok {
		return f.Name
	}
	opts := strings.Split(val, ",")
	fieldTag = opts[0]
	if fieldTag == "-" {
		return ""
	}
	return fieldTag
}

//*Struct or(Struct) ==>map[tagName]interface{} like map[db_name]interface{} ;tag like db
func TransFormByTagToMap(ptrs interface{}, tag string) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	t := reflect.TypeOf(ptrs)
	v := reflect.ValueOf(ptrs)

	if v.IsNil() {
		return ret, nil
	}
	if t.Elem().Kind() != reflect.Struct {
		return ret, errors.New("ptrs must struct ptr")
	}
	retVal := reflect.ValueOf(ret)
	num := v.Elem().NumField()

	for i := 0; i < num; i++ {

		if v.Elem().Field(i).Type().Kind() == reflect.Ptr {
			if k := getTagName(t.Elem().Field(i), tag); k != "" && k != "-" {
				v := v.Elem().Field(i)
				if v.Type().Kind() == reflect.Ptr && !v.IsNil() {
					if !v.IsNil() {
						retVal.SetMapIndex(reflect.ValueOf(k), v.Elem())
					}
				}
			}
		} else {
			if k := getTagName(t.Elem().Field(i), tag); k != "" && k != "-" {
				retVal.SetMapIndex(reflect.ValueOf(k), v.Elem().Field(i))
			}
		}

	}

	return ret, nil
}

/*
dbString==>Struct or *Struct;
eg:
 - args  in  like {"id":"1","name":"zs","h":"1.12"}; key=db_field_name;
 - args ptrs like &Struct{};
*/
func ScanDbDataToStruct(in string, ptrs interface{}) error {
	inj := map[string]string{}
	err := json.Unmarshal([]byte(in), &inj)
	if err != nil {
		return err
	}

	t := reflect.TypeOf(ptrs)
	v := reflect.ValueOf(ptrs)

	if v.IsNil() {
		return nil
	}
	if t.Elem().Kind() != reflect.Struct {
		return errors.New("ptrs must struct ptr")
	}

	num := v.Elem().NumField()
	for i := 0; i < num; i++ {
		if k := getTagName(t.Elem().Field(i), "db"); k != "" {
			if "" == inj[k] {
				continue
			}
			if err := changeType(inj[k], v.Elem().Field(i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func changeType(s string, dv reflect.Value) error {
	flag := false
	tmp := dv
	if tmp.Kind() == reflect.Ptr {
		flag = true
		tmp = getKind(dv)
	}

	switch tmp.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		i64, err := strconv.ParseInt(s, 10, tmp.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", "string", s, dv.Kind(), err)
		}
		if flag {
			dv.Set(reflect.ValueOf(&i64))
		} else {
			dv.SetInt(i64)
		}

		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		u64, err := strconv.ParseUint(s, 10, tmp.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", "string", s, dv.Kind(), err)
		}

		if flag {
			dv.Set(reflect.ValueOf(&u64))
		} else {
			dv.SetUint(u64)
		}
		return nil
	case reflect.Float32, reflect.Float64:

		f64, err := strconv.ParseFloat(s, tmp.Type().Bits())
		if err != nil {
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", "string", s, dv.Kind(), err)
		}

		if flag {
			dv.Set(reflect.ValueOf(&f64))
		} else {
			dv.SetFloat(f64)
		}
		return nil
	case reflect.String:
		if flag {
			dv.Set(reflect.ValueOf(&s))
		} else {
			dv.SetString(s)
		}
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", "string", dv.Type())
}

func getKind(value reflect.Value) reflect.Value {
	switch value.Type().String() {
	case "*int", "*int8", "*int16", "*int32", "*int64":
		return reflect.ValueOf(int64(DEFAULT_NUM))

	case "*uint", "*uint8", "*uint16", "*uint32", "*uint64":
		return reflect.ValueOf(uint64(DEFAULT_NUM))
	case "*float32", "*float64":
		return reflect.ValueOf(float64(DEFAULT_NUM))
	case "*string":
		return reflect.ValueOf(DEFAULT_STR)
	default:
		return reflect.ValueOf(DEFAULT_STR)
	}
}
