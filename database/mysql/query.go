package mysql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

//反射一个结构体的成员指针，包括非导出成员
func StructFieldsPointerOf(obj interface{}) ([]interface{}, error) {
	if obj == nil {
		return nil, errors.New("nil value")
	}
	v := reflect.ValueOf(obj)
	t := v.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, errors.New("not pointer or not struct pointer")
	}
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	vals := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		pv := v.Field(i)
		vals[i] = reflect.NewAt(pv.Type(), unsafe.Pointer(pv.UnsafeAddr())).Interface()
	}
	return vals, nil
}

type Query interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

//从一个可查询mysql数据的query中将数据写入o
func SearchFrom(o interface{}, query Query, sql string, args ...interface{}) error {
	tp := reflect.TypeOf(o)
	if tp.Kind() != reflect.Ptr {
		return errors.New("o is not a pointer")
	}

	if tp.Elem().Kind() == reflect.Slice &&
		tp.Elem().Elem().Kind() == reflect.Ptr &&
		tp.Elem().Elem().Elem().Kind() == reflect.Struct { //[]*struct
		t := tp.Elem().Elem().Elem()
		rows, err := query.Query(sql, args...)
		if err != nil {
			return err
		}
		src := reflect.Indirect(reflect.ValueOf(o))
		for rows.Next() {
			v := reflect.New(t)
			p, _ := StructFieldsPointerOf(v.Interface())
			err = rows.Scan(p...)
			if err != nil {
				return err
			}
			src.Set(reflect.Append(src, v))
		}
	} else if tp.Elem().Kind() == reflect.Struct {
		row := query.QueryRow(sql, args...)
		p, _ := StructFieldsPointerOf(o)
		if err := row.Scan(p...); err != nil {
			return err
		}
	} else {

		switch tp.Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String:
			row := query.QueryRow(sql, args...)
			if err := row.Scan(o); err != nil {
				return err
			}
			return nil
		}
		
		return errors.New("not support type of o")
	}
	return nil
}


//mysql visible
//使得sql语句结构化
func NewSql() *Sql {
	return &Sql{}
}

type Sql struct {
	sl, from, where, groupBy, orderBy, limit string
}

func (sl *Sql) Select(s string) *Sql {
	sl.sl = s
	return sl
}
func (sl *Sql) SelectSlice(s []string) *Sql {
	sl.sl = strings.Join(s, ",")
	return sl
}

func (sl *Sql) From(s string) *Sql {
	sl.from = s
	return sl
}
func (sl *Sql) Where(s string) *Sql {
	sl.where = s
	return sl
}
func (sl *Sql) And(s string) *Sql {
	sl.where += " AND " + s
	return sl
}
func (sl *Sql) GroupBy(s string) *Sql {
	sl.groupBy = s
	return sl
}
func (sl *Sql) OrderBy(s string) *Sql {
	sl.orderBy = s
	return sl
}
func (sl *Sql) Limit(s string) *Sql {
	sl.limit = s
	return sl
}
func (sl *Sql) Sql() string {
	str := bytes.Buffer{}
	str.WriteString(fmt.Sprintf("select %s from %s ", sl.sl, sl.from))
	if sl.where != "" {
		str.WriteString(" where ")
		str.WriteString(sl.where)
	}
	if sl.groupBy != "" {
		str.WriteString(" group by  ")
		str.WriteString(sl.groupBy)
	}
	if sl.orderBy != "" {
		str.WriteString(" order by  ")
		str.WriteString(sl.orderBy)
	}
	if sl.limit != "" {
		str.WriteString(" limit  ")
		str.WriteString(sl.limit)
	}
	return str.String()
}

func (sl *Sql) ResultIn(o interface{}, args ...interface{}) error {
	return SearchFrom(o, GetOneUsableDb(), sl.Sql(), args...)
}

func (sl *Sql) ResultInDb(o interface{}, query Query, args ...interface{}) error {
	return SearchFrom(o, query, sl.Sql(), args...)
}
