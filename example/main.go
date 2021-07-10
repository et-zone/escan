package main

import (
	"database/sql"
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
	// Test_Insert()
	// Test_Select()
	Test_Update()
	// Test_Del()
}

func Test_Insert() {
	s := []interface{}{}
	name := "biubiu"
	age := float64(12)
	t := time.Now().Format("2006-01-02 15:04:05")
	tmp := Stu{
		Name:  name,
		Age:   age,
		Ctime: t,
	}
	s = append(s, tmp)
	build := escan.NewBuilder("stu", new(Stu))
	sql, args := build.InsertBuilderSql(&s)
	_, err := sqlDB.Exec(sql, args...)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func Test_Update() {

	kv := map[string]interface{}{"age": 18}
	var id interface{} = escan.Int64(5)
	condtions := map[string]escan.Condition{"id": escan.Condition{Equal: &id}}

	build := escan.NewBuilder("stu", new(Stu))
	sql, args := build.UpdateBuilderSql(kv, condtions)
	_, err := sqlDB.Exec(sql, args...)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func Test_Select() {
	fields := []string{"name"} //写了要保证准确
	// fields := []string{} //不写默认全部查询
	// var id interface{} = 6
	// condtions := map[string]db.Condition{"id": db.Condition{Equal: &id}}
	condtions := map[string]escan.Condition{}
	page := 2
	pageSize := 3
	sr := &escan.Screen{
		Limit:       pageSize,
		OfSet:       (page - 1) * pageSize,
		OrderByDesc: []string{"id"},
	}
	build := escan.NewBuilder("stu", new(Stu))
	sql, args := build.SelectBuilderSql(fields, condtions, sr)

	rows, err := sqlDB.Query(sql, args...)
	if err != nil {
		fmt.Println(err.Error())
	}
	stus := []Stu{}
	err = escan.NewEscan().ScanAll(&stus, rows)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(stus)
}

func Test_Del() {

	var id interface{} = 6
	condtions := map[string]escan.Condition{"id": escan.Condition{Equal: &id}}

	build := escan.NewBuilder("stu", new(Stu))
	sql, args := build.DeleteBuilderSql(condtions)
	_, err := sqlDB.Exec(sql, args...)
	if err != nil {
		fmt.Println(err.Error())
	}

}
