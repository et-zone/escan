package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.con/et-zone/escan"
)

var sqlDB *sql.DB

type EConfig struct {
	UserName string `json:"userName"`
	PassWord string `json:"passWord"`
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	DB       string `json:"db"`
}

type Stu struct {
	ID    int64   `json:"id" db:"id" fieldtag:"select"` //db=- 不会查询,不带tag不操作。db 和fieldtag 必传
	Name  string  `json:"name" db:"name" fieldtag:"insert,select"`
	Age   float64 `json:"age" db:"age" fieldtag:"insert,select"`
	Ctime string  `json:"ctime" db:"c_time" fieldtag:"insert,select"` //
}

//结构化对外可以使用，对外是json的就不需要使用了
type StuType struct {
	ID    *int64   `json:"id" db:"id" fieldtag:"select"` //db=- 不会查询,不带tag不操作。db 和fieldtag 必传
	Name  *string  `json:"name" db:"name" fieldtag:"insert,select"`
	Age   *float64 `json:"age" db:"age" fieldtag:"insert,select"`
	Ctime *string  `json:"ctime" db:"c_time" fieldtag:"insert,select"` //
}

var stuBuild = escan.NewBuilder("stu", new(Stu))

func InsertStus(des *[]interface{}) error {
	sql, args := stuBuild.InsertBuilderSql(des)
	_, err := sqlDB.Exec(sql, args...)
	return err
}

func UpdateStu(kv map[string]interface{}, conditions map[string]*escan.Condition) error {
	sql, args := stuBuild.UpdateBuilderSql(kv, conditions)
	_, err := sqlDB.Exec(sql, args...)
	return err
}

func DelStu(conditions map[string]*escan.Condition) error {
	sql, args := stuBuild.DeleteBuilderSql(conditions)
	_, err := sqlDB.Exec(sql, args...)
	return err
}

func SelectStus(conditions map[string]*escan.Condition, screen *escan.Screen) (*[]Stu, error) {

	sql, args := stuBuild.SelectBuilderSql([]string{}, conditions, screen)
	rows, err := sqlDB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	stus := []Stu{}
	err = escan.NewEscan().ScanAll(&stus, rows)
	if err != nil {
		return nil, err
	}
	return &stus, err
}

//fields Can be empty,If it is empty, it will get all
func SelectStusToMap(fields []string, conditions map[string]*escan.Condition, screen *escan.Screen) (*[]map[string]string, error) {

	sql, args := stuBuild.SelectBuilderSql(fields, conditions, screen)
	rows, err := sqlDB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	stusMap := []map[string]string{}
	err = escan.NewEscan().ScanAll(&stusMap, rows)
	if err != nil {
		return nil, err
	}
	for i, val := range stusMap {
		v, err := escan.ChToJsonByTagDB(val, Stu{})
		if err != nil {
			stusMap[i] = nil
		} else {
			stusMap[i] = v
		}
	}
	return &stusMap, err
}

//field Can be '',If it is '', it will get count(*)
func SelectStusCount(field string, conditions map[string]escan.Condition, screen *escan.Screen) (int, error) {

	sql, args := stuBuild.SelectBuilderCountSql(field, conditions, screen)

	rows, err := sqlDB.Query(sql, args...)
	if err != nil {
		return 0, err
	}
	count := map[string]int{}
	err = escan.NewEscan().ScanOne(&count, rows)
	if err != nil {
		return 0, err
	}
	return count["c"], err
}

func InitSQL(cfg *EConfig) {
	sqlClient, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			cfg.UserName,
			cfg.PassWord,
			cfg.Addr,
			cfg.Port,
			cfg.DB,
		),
	)
	if err != nil {
		panic(err)
	}
	client := sqlClient
	client.SetMaxIdleConns(30)
	client.SetMaxOpenConns(100)
	fmt.Println("init sql success")

	sqlDB = client

}

func main() {
	InitSQL(&EConfig{
		UserName: "root",
		PassWord: "mysql",
		Addr:     "49.232.190.114",
		Port:     3366,
		DB:       "test",
	})
	Test_Insert()

	// Test_Update()
	// Test_Del()

	Test_Select_Map()
	// Test_Select_Stu()
	// Test_Select_Count()
	// s := "2006-01-02 15:04:05"
	// fmt.Println(s[11:])
}

func Test_Insert() {
	s := []interface{}{}
	name := "biubiu99"
	// age := float64(12)
	t := time.Now().Format("2006-01-02 15:04:05")
	tmp := Stu{
		Name: name,
		// Age:   age,
		Ctime: t,
	}

	// tmpType := StuType{
	// 	Name: &name,
	// 	Age:  &age,
	// }

	s = append(s, tmp)

	err := InsertStus(&s)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func Test_Update() {

	kv := map[string]interface{}{"age": 18}
	// var id interface{} = escan.Int64(22)
	conditions := map[string]*escan.Condition{"id": escan.NewCondition().Equal(18)}

	err := UpdateStu(kv, conditions)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func Test_Select_Map() {
	// fields := []string{"id"} //写了要保证准确
	fields := []string{} //不写默认全部查询
	// var id interface{} = 6
	// condtions := map[string]db.Condition{"id": db.Condition{Equal: &id}}
	conditions := map[string]*escan.Condition{}
	page := 2
	pageSize := 3
	sr := &escan.Screen{}
	sr.SetOrderByDesc([]string{"id"}).SetPageSize(page, pageSize)
	data, err := SelectStusToMap(fields, conditions, sr)
	if err != nil {
		fmt.Println(err.Error())
	}
	b, _ := json.Marshal(data)
	fmt.Println(string(b))
}

func Test_Select_Stu() {
	conditions := map[string]*escan.Condition{}
	page := 2
	pageSize := 3
	sr := &escan.Screen{}
	sr.SetOrderByDesc([]string{"id"}).SetPageSize(page, pageSize)
	data, err := SelectStus(conditions, sr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(*data)
}

func Test_Del() {

	// var id interface{} = 2
	conditions := map[string]*escan.Condition{
		"id": escan.NewCondition().Equal(18),
	}
	b, _ := json.Marshal(escan.NewCondition().Equal(18))
	fmt.Println(string(b))
	err := DelStu(conditions)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func Test_Select_Count() {
	conditions := map[string]escan.Condition{}
	page := 2
	pageSize := 3
	sr := &escan.Screen{}
	sr.SetOrderByDesc([]string{"id"}).SetPageSize(page, pageSize)
	data, err := SelectStusCount("", conditions, sr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(data)
}
