package dba

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"teampilot/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type QueryBuilder struct {
	*gorm.DB
}

func GetModelEntiryFields(model interface{}) string {
	s, _ := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	fields := []string{}
	for _, v := range s.Fields {
		if v.DBName != "" {
			fields = append(fields, fmt.Sprintf(`"%s".%s`, s.Table, v.DBName))
		}
	}

	return strings.Join(fields, ",")
}

// Total is return total from query
func (db *QueryBuilder) Total() int {
	var total int64
	db.Count(&total)
	return int(total)
}

func (db *QueryBuilder) Sum(fieldName string) float64 {
	var total float64
	db.Select(fmt.Sprintf("sum(%s)", fieldName)).Row().Scan(&total)
	return total
}

func (db *QueryBuilder) Like(field string, value *string) *QueryBuilder {
	if value != nil {
		db.DB = db.Where(fmt.Sprintf("LOWER(%s) like LOWER(?)", field), "%"+*value+"%")
	}
	return db
}

func (db *QueryBuilder) NotNull(field string) *QueryBuilder {
	db.DB = db.Where(fmt.Sprintf("%s notnull", field))
	return db
}

func (db *QueryBuilder) Null(field string) *QueryBuilder {
	db.DB = db.Where(fmt.Sprintf("%s is null", field))
	return db
}

func (db *QueryBuilder) CastLike(field string, value *string) *QueryBuilder {
	if value != nil {
		db.DB = db.Where(fmt.Sprintf("CAST(%s as TEXT) like ?", field), "%"+*value+"%")
	}
	return db
}

func (db *QueryBuilder) OrLike(field string, value *string) *QueryBuilder {
	if value != nil {
		db.DB = db.Or(fmt.Sprintf("LOWER(%s) like LOWER(?)", field), "%"+*value+"%")
	}
	return db
}

func (db *QueryBuilder) Equal(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s = ?", field), value)
	}
	return db
}

func (db *QueryBuilder) DateEqual(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s::DATE = ?", field), value)
	}
	return db
}

func (db *QueryBuilder) DateLowerEqual(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s::DATE <= ?", field), value)
	}
	return db
}

func (db *QueryBuilder) YearMonthEqual(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("to_char(%s, 'YYYY-MM') = substring(?,0,8)", field), value)
	}
	return db
}

func (db *QueryBuilder) NotEqual(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s <> ?", field), value)
	}
	return db
}

func (db *QueryBuilder) Between(field string, value interface{}) *QueryBuilder {
	s := reflect.ValueOf(value)
	if s.Kind() != reflect.Slice {
		return db
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return db
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	db.DB = db.Where(fmt.Sprintf("%s between ? and ?", field), ret[0], ret[1])
	return db
}

func (db *QueryBuilder) BetweenDate(field string, value interface{}) *QueryBuilder {
	s := reflect.ValueOf(value)
	if s.Kind() != reflect.Slice {
		return db
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return db
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}
	db.DB = db.Where(fmt.Sprintf("%s::DATE between ? and ?", field), ret[0], ret[1])
	return db
}

func (db *QueryBuilder) BetweenTime(field string, value interface{}) *QueryBuilder {
	s := reflect.ValueOf(value)
	if s.Kind() != reflect.Slice {
		return db
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return db
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	db.DB = db.Where(fmt.Sprintf("date_part('hour', %s) >= ? and date_part('hour', %s) <= ?", field, field), ret[0], ret[1])
	return db
}

// value can be empty
func (db *QueryBuilder) Bool(field string, value *string) *QueryBuilder {
	if !utils.IsNil(value) {
		if *value == "1" {
			db.DB = db.Where(fmt.Sprintf("%s = ?", field), true)
		} else if *value == "0" {
			db.DB = db.Where(fmt.Sprintf("%s = ?", field), false)
		}
	}
	return db
}

func (db *QueryBuilder) Include(field string, values interface{}) *QueryBuilder {
	if !utils.IsNil(values) {
		db.DB = db.Where(fmt.Sprintf("%s IN ?", field), values)
	}
	return db
}

func (db *QueryBuilder) IncludeSubquery(field string, values interface{}) *QueryBuilder {
	if !utils.IsNil(values) {
		db.DB = db.Where(fmt.Sprintf("%s IN (?)", field), values)
	}
	return db
}

func (db *QueryBuilder) NotInclude(field string, values interface{}) *QueryBuilder {
	if !utils.IsNil(values) {
		db.DB = db.Where(fmt.Sprintf("%s NOT IN ?", field), values)
	}
	return db
}

func (db *QueryBuilder) OrEqual(field string, value interface{}) *QueryBuilder {
	if !(value == nil || reflect.ValueOf(value).IsNil()) {
		db.DB = db.Or(fmt.Sprintf("%s = ?", field), value)
	}
	return db
}

func (db *QueryBuilder) Greater(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s > ?", field), value)
	}
	return db
}

func (db *QueryBuilder) EqualGreater(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s >= ?", field), value)
	}
	return db
}

func (db *QueryBuilder) EqualInPqArray(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("? = ANY (%s)", field), value)
	}
	return db
}

func (db *QueryBuilder) OrEqualGreater(field string, value interface{}) *QueryBuilder {
	if !(value == nil || reflect.ValueOf(value).IsNil()) {
		db.DB = db.Or(fmt.Sprintf("%s >= ?", field), value)
	}
	return db
}

func (db *QueryBuilder) EqualLower(field string, value interface{}) *QueryBuilder {
	if !utils.IsNil(value) {
		db.DB = db.Where(fmt.Sprintf("%s <= ?", field), value)
	}
	return db
}

func (db *QueryBuilder) OrEqualLower(field string, value interface{}) *QueryBuilder {
	if !(value == nil || reflect.ValueOf(value).IsNil()) {
		db.DB = db.Or(fmt.Sprintf("%s <= ?", field), value)
	}
	return db
}

func (db *QueryBuilder) Search(fields []string, value *string) *QueryBuilder {
	if !utils.IsNil(value) {
		queryString := ""
		for index, field := range fields {
			if index == 0 {
				queryString = fmt.Sprintf("LOWER(%s) like LOWER(@value)", field)
				continue
			}
			fieldQuery := fmt.Sprintf("LOWER(%s) like LOWER(@value)", field)
			queryString = fmt.Sprintf("%s or %s", queryString, fieldQuery)
		}
		db.DB = db.Where(queryString, sql.Named("value", "%"+*value+"%"))
	}

	return db
}
