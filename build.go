package escan

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
)

const (
	TAG_INSERT = "insert"
	TAG_SELECT = "select"
)

//update value
type KV map[string]interface{}

//where condition
type Condition struct {
	//equal
	Eq *interface{} `json:"equal,omitempty"`
	//contain
	Ct []interface{} `json:"contain,omitempty"`
	//notEqual
	Ne *interface{} `json:"neq,omitempty"`
	//lessThan
	Lt *interface{} `json:"lt,omitempty"`
	//lessEqualThan
	Let *interface{} `json:"lte,omitempty"`
	//GreaterThan
	Gt *interface{} `json:"gt,omitempty"`
	//GreaterEqualThan
	Get *interface{} `json:"gte,omitempty"`
	//between
	Bet []interface{} `json:"between,omitempty"`
	//like
	Lk *string `json:"like,omitempty"`
}

func NewCondition() *Condition {
	return &Condition{}
}

func (cd *Condition) Equal(e interface{}) *Condition {
	cd.Eq = &e
	return cd
}

func (cd *Condition) NotEqual(neq interface{}) *Condition {
	cd.Ne = &neq
	return cd
}

func (cd *Condition) Contain(ct []interface{}) *Condition {
	cd.Ct = ct
	return cd
}

func (cd *Condition) LessThan(lt interface{}) *Condition {
	cd.Lt = &lt
	return cd
}

func (cd *Condition) LessEqualThan(lte interface{}) *Condition {
	cd.Let = &lte
	return cd
}

func (cd *Condition) GreaterThan(gt interface{}) *Condition {
	cd.Gt = &gt
	return cd
}

func (cd *Condition) GreaterEqualThan(lte interface{}) *Condition {
	cd.Get = &lte
	return cd
}

func (cd *Condition) Between(between []interface{}) *Condition {
	cd.Bet = between
	return cd
}

func (cd *Condition) Like(like string) *Condition {
	cd.Lk = &like
	return cd
}

//Screen condition
type Screen struct {
	OrderByAsc  []string `json:"orderByAsc,omitempty"`
	OrderByDesc []string `json:"orderByDesc,omitempty"`
	Limit       int      `json:"like,omitempty"`
	OfSet       int      `json:"like,omitempty"`
}

func NewScreen() Screen {
	return Screen{}
}

//like page=1,size=10
func (this *Screen) SetPageSize(page int, pageSize int) *Screen {
	this.OfSet = int((page - 1) * pageSize)
	this.Limit = int(pageSize)
	return this
}

func (this *Screen) SetOrderByAsc(fields []string) *Screen {
	this.OrderByAsc = fields
	return this
}
func (this *Screen) SetOrderByDesc(fields []string) *Screen {
	this.OrderByDesc = fields
	return this
}

type SBuilder struct {
	TableName string
	st        *sqlbuilder.Struct
}

func NewBuilder(tableName string, structValue interface{}) *SBuilder {
	return &SBuilder{
		TableName: tableName,
		st:        sqlbuilder.NewStruct(structValue),
	}

}

func (this *SBuilder) UpdateBuilderSql(kv KV, conditions map[string]*Condition) (sql string, args []interface{}) {
	up := sqlbuilder.Update(this.TableName)
	assigns := []string{}
	for k, v := range kv {
		if v != nil {
			assigns = append(assigns, up.Assign(k, v))
		}
	}
	up.Set(assigns...)
	for field, condition := range conditions {
		if len(condition.Ct) != 0 {
			up.Where(up.In(field, condition.Ct...))
		}
		if condition.Eq != nil {
			up.Where(up.Equal(field, *condition.Eq))
		}
		if condition.Ne != nil {
			up.Where(up.NotEqual(field, *condition.Ne))
		}
		if len(condition.Bet) == 2 {
			up.Where(up.Between(field, condition.Bet[0], condition.Bet[1]))
		}
		if condition.Lk != nil {
			up.Where(up.Like(field, *condition.Lk))
		}
		if condition.Lt != nil {
			up.Where(up.LessThan(field, *condition.Lt))
		}
		if condition.Let != nil {
			up.Where(up.LessEqualThan(field, *condition.Let))
		}
		if condition.Gt != nil {
			up.Where(up.GreaterThan(field, *condition.Gt))
		}
		if condition.Get != nil {
			up.Where(up.GreaterEqualThan(field, *condition.Get))
		}

	}

	return up.Build()
}

func (this *SBuilder) SelectBuilderSql(fields []string, conditions map[string]*Condition, screen *Screen) (sql string, args []interface{}) {
	se := this.st.SelectFromForTag(this.TableName, TAG_SELECT)
	if len(fields) != 0 {
		se.Select(fields...)
	}
	for field, condition := range conditions {
		if len(condition.Ct) != 0 {
			se.Where(se.In(field, condition.Ct...))
		}
		if condition.Eq != nil {
			se.Where(se.Equal(field, *condition.Eq))
		}
		if condition.Ne != nil {
			se.Where(se.NotEqual(field, *condition.Ne))
		}
		if len(condition.Bet) == 2 {
			se.Where(se.Between(field, condition.Bet[0], condition.Bet[1]))
		}
		if condition.Lk != nil {
			se.Where(se.Like(field, *condition.Lk))
		}
		if condition.Lt != nil {
			se.Where(se.LessThan(field, *condition.Lt))
		}
		if condition.Let != nil {
			se.Where(se.LessEqualThan(field, *condition.Let))
		}
		if condition.Gt != nil {
			se.Where(se.GreaterThan(field, *condition.Gt))
		}
		if condition.Get != nil {
			se.Where(se.GreaterEqualThan(field, *condition.Get))
		}
	}
	if screen != nil {
		if len(screen.OrderByAsc) > 0 {
			se.OrderBy(screen.OrderByAsc...).Asc()
		}
		if len(screen.OrderByDesc) > 0 {
			se.OrderBy(screen.OrderByDesc...).Desc()
		}
		if screen.Limit > 0 {
			se.Limit(screen.Limit)
		}
		if screen.Limit > 0 {
			se.Offset(screen.OfSet)
		}
	}
	return se.Build()
}

func (this *SBuilder) InsertBuilderSql(des *[]interface{}) (sql string, args []interface{}) {
	return this.st.InsertIntoForTag(this.TableName, TAG_INSERT, *des...).Build()
}

func (this *SBuilder) DeleteBuilderSql(conditions map[string]*Condition) (sql string, args []interface{}) {
	del := sqlbuilder.DeleteFrom(this.TableName)

	for field, condition := range conditions {
		if len(condition.Ct) != 0 {
			del.Where(del.In(field, condition.Ct...))
		}
		if condition.Eq != nil {
			del.Where(del.Equal(field, *condition.Eq))
		}
		if condition.Ne != nil {
			del.Where(del.NotEqual(field, *condition.Ne))
		}
		if len(condition.Bet) == 2 {
			del.Where(del.Between(field, condition.Bet[0], condition.Bet[1]))
		}
		if condition.Lk != nil {
			del.Where(del.Like(field, *condition.Lk))
		}
		if condition.Lt != nil {
			del.Where(del.LessThan(field, *condition.Lt))
		}
		if condition.Let != nil {
			del.Where(del.LessEqualThan(field, *condition.Let))
		}
		if condition.Gt != nil {
			del.Where(del.GreaterThan(field, *condition.Gt))
		}
		if condition.Get != nil {
			del.Where(del.GreaterEqualThan(field, *condition.Get))
		}

	}
	return del.Build()
}

//if field is "" will be select count(*)
func (this *SBuilder) SelectBuilderCountSql(field string, conditions map[string]Condition, screen *Screen) (sql string, args []interface{}) {
	se := this.st.SelectFromForTag(this.TableName, TAG_SELECT)
	if field == "" {
		se.Select("count(*) as c")
	} else {
		se.Select(fmt.Sprintf("count(%s) as c", field))
	}

	for field, condition := range conditions {
		if len(condition.Ct) != 0 {
			se.Where(se.In(field, condition.Ct...))
		}
		if condition.Eq != nil {
			se.Where(se.Equal(field, *condition.Eq))
		}
		if condition.Ne != nil {
			se.Where(se.NotEqual(field, *condition.Ne))
		}
		if len(condition.Bet) == 2 {
			se.Where(se.Between(field, condition.Bet[0], condition.Bet[1]))
		}
		if condition.Lk != nil {
			se.Where(se.Like(field, *condition.Lk))
		}
		if condition.Lt != nil {
			se.Where(se.LessThan(field, *condition.Lt))
		}
		if condition.Let != nil {
			se.Where(se.LessEqualThan(field, *condition.Let))
		}
		if condition.Gt != nil {
			se.Where(se.GreaterThan(field, *condition.Gt))
		}
		if condition.Get != nil {
			se.Where(se.GreaterEqualThan(field, *condition.Get))
		}
	}
	if screen != nil {
		if len(screen.OrderByAsc) > 0 {
			se.OrderBy(screen.OrderByAsc...).Asc()
		}
		if len(screen.OrderByDesc) > 0 {
			se.OrderBy(screen.OrderByDesc...).Desc()
		}

	}
	return se.Build()
}
