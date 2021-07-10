package escan

import (
	"github.com/huandu/go-Sqlbuilder"
)

const (
	TAG_INSERT = "insert"
	TAG_SELECT = "select"
)

//update value
type KV map[string]interface{}

//where condition
type Condition struct {
	Equal            *interface{}  `json:"equal,omitempty"`
	Contain          []interface{} `json:"contain,omitempty"`
	NotEqual         *interface{}  `json:"neq,omitempty"`
	LessThan         *interface{}  `json:"lt,omitempty"`
	LessEqualThan    *interface{}  `json:"lte,omitempty"`
	GreaterThan      *interface{}  `json:"gt,omitempty"`
	GreaterEqualThan *interface{}  `json:"gte,omitempty"`
	Between          []interface{} `json:"between,omitempty"`
	Like             *string       `json:"like,omitempty"`
}

//Screen condition
type Screen struct {
	OrderByAsc  []string `json:"orderByAsc,omitempty"`
	OrderByDesc []string `json:"orderByDesc,omitempty"`
	Limit       int      `json:"like,omitempty"`
	OfSet       int      `json:"like,omitempty"`
}

//like page=1,size=10
func (this *Screen) SetPageSize(page int64, pageSize int64) {
	this.OfSet = int((page - 1) * pageSize)
	this.Limit = int(pageSize)
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

func (this *SBuilder) UpdateBuilderSql(kv KV, conditions map[string]Condition) (sql string, args []interface{}) {
	up := sqlbuilder.Update(this.TableName)
	assigns := []string{}
	for k, v := range kv {
		if v != nil {
			assigns = append(assigns, up.Assign(k, v))
		}
	}
	up.Set(assigns...)
	for field, condition := range conditions {
		if len(condition.Contain) != 0 {
			up.Where(up.In(field, condition.Contain...))
		}
		if condition.Equal != nil {
			up.Where(up.Equal(field, *condition.Equal))
		}
		if condition.NotEqual != nil {
			up.Where(up.NotEqual(field, *condition.NotEqual))
		}
		if len(condition.Between) == 2 {
			up.Where(up.Between(field, condition.Between[0], condition.Between[1]))
		}
		if condition.Like != nil {
			up.Where(up.Like(field, *condition.Like))
		}
		if condition.LessThan != nil {
			up.Where(up.LessThan(field, *condition.LessThan))
		}
		if condition.LessEqualThan != nil {
			up.Where(up.LessEqualThan(field, *condition.LessEqualThan))
		}
		if condition.GreaterThan != nil {
			up.Where(up.GreaterThan(field, *condition.GreaterThan))
		}
		if condition.GreaterEqualThan != nil {
			up.Where(up.GreaterEqualThan(field, *condition.GreaterEqualThan))
		}

	}

	return up.Build()
}

func (this *SBuilder) SelectBuilderSql(fields []string, conditions map[string]Condition, screen *Screen) (sql string, args []interface{}) {
	se := this.st.SelectFromForTag(this.TableName, TAG_SELECT)
	if len(fields) != 0 {
		se.Select(fields...)
	}
	for field, condition := range conditions {
		if len(condition.Contain) != 0 {
			se.Where(se.In(field, condition.Contain...))
		}
		if condition.Equal != nil {
			se.Where(se.Equal(field, *condition.Equal))
		}
		if condition.NotEqual != nil {
			se.Where(se.NotEqual(field, *condition.NotEqual))
		}
		if len(condition.Between) == 2 {
			se.Where(se.Between(field, condition.Between[0], condition.Between[1]))
		}
		if condition.Like != nil {
			se.Where(se.Like(field, *condition.Like))
		}
		if condition.LessThan != nil {
			se.Where(se.LessThan(field, *condition.LessThan))
		}
		if condition.LessEqualThan != nil {
			se.Where(se.LessEqualThan(field, *condition.LessEqualThan))
		}
		if condition.GreaterThan != nil {
			se.Where(se.GreaterThan(field, *condition.GreaterThan))
		}
		if condition.GreaterEqualThan != nil {
			se.Where(se.GreaterEqualThan(field, *condition.GreaterEqualThan))
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

func (this *SBuilder) DeleteBuilderSql(conditions map[string]Condition) (sql string, args []interface{}) {
	del := sqlbuilder.DeleteFrom(this.TableName)

	for field, condition := range conditions {
		if len(condition.Contain) != 0 {
			del.Where(del.In(field, condition.Contain...))
		}
		if condition.Equal != nil {
			del.Where(del.Equal(field, *condition.Equal))
		}
		if condition.NotEqual != nil {
			del.Where(del.NotEqual(field, *condition.NotEqual))
		}
		if len(condition.Between) == 2 {
			del.Where(del.Between(field, condition.Between[0], condition.Between[1]))
		}
		if condition.Like != nil {
			del.Where(del.Like(field, *condition.Like))
		}
		if condition.LessThan != nil {
			del.Where(del.LessThan(field, *condition.LessThan))
		}
		if condition.LessEqualThan != nil {
			del.Where(del.LessEqualThan(field, *condition.LessEqualThan))
		}
		if condition.GreaterThan != nil {
			del.Where(del.GreaterThan(field, *condition.GreaterThan))
		}
		if condition.GreaterEqualThan != nil {
			del.Where(del.GreaterEqualThan(field, *condition.GreaterEqualThan))
		}

	}
	return del.Build()
}
