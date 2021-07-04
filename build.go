package escan

import (
	"github.com/huandu/go-Sqlbuilder"
)

const (
	TAG_INSERT = "insert"
	TAG_SELECT = "select"
)

func NewSelectBuilder(tableName string, tag string, up *sqlbuilder.SelectBuilder) {

}

func NewUpdateBuilder(tableName string, tag string, up *sqlbuilder.SelectBuilder) {

}

//update value
type KV map[string]interface{}

//where condition
type Condition struct {
	Equal            *interface{}
	Contain          []interface{}
	NotEqual         *interface{}
	LessThan         *interface{}
	LessEqualThan    *interface{}
	GreaterThan      *interface{}
	GreaterEqualThan *interface{}
	Between          []interface{}
	Like             *string
}

//Screen condition
type Screen struct {
	OrderByAsc  []string
	OrderByDesc []string
	Asc         bool
	Limit       int
	OfSet       int
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
