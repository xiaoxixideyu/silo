package search

import (
	"reflect"
	"silo/pkg/utils"
	"strings"
)

// Tag .
type Tag struct {
	Table          string
	As             string
	Operator       string
	Column         string
	JsonValue      string
	On             [][]string // join on condition
	OnRaw          string
	Join           string // join target table name
	Value          any
	Children       []*Tag
	SkipCheckEmpty bool // skip check empty value
}

// AddChild .
func (t *Tag) AddChild(tag *Tag) *Tag {
	if t.Children == nil {
		t.Children = make([]*Tag, 0)
	}

	t.Children = append(t.Children, tag)

	return t
}

// ParseSearchTag .
func ParseSearchTag(tableName string, q any, condition Condition, parent ...*Tag) {
	valueOf := reflect.ValueOf(q)

	var qType reflect.Type
	var qValue reflect.Value

	// handler pointer or struct
	if valueOf.Kind() == reflect.Pointer {
		qType = reflect.TypeOf(q).Elem()
		qValue = reflect.ValueOf(q).Elem()
	} else if valueOf.Kind() == reflect.Struct {
		qType = reflect.TypeOf(q)
		qValue = reflect.ValueOf(q)
	}

	var tag string
	var ok bool
	if qValue.Kind() == reflect.Struct {
		for i := 0; i < qType.NumField(); i++ {
			typeField := qType.Field(i)
			valueField := qValue.Field(i)
			// skip unexported field
			if !typeField.IsExported() {
				continue
			}
			tag, ok = "", false
			tag, ok = typeField.Tag.Lookup(QueryTag)
			if !ok {
				// parse child field
				ParseSearchTag(tableName, valueField.Interface(), condition)
				continue
			}
			switch tag {
			case "-":
				continue
			}
			t := makeTag(tableName, tag)

			isStruct := utils.IsStruct(valueField)

			// join之类的操作还是需要将当前struct定义加入tag进行处理
			// skip zero field if its not a struct
			if !t.SkipCheckEmpty && !isStruct && valueField.IsZero() {
				continue
			}

			if len(parent) == 0 {
				condition.AddSearchTag(t)
			} else {
				parent[0].AddChild(t)
			}

			switch t.Operator {
			case OPLeftJoin, OPJoin, OPGroupOr:
				ParseSearchTag(tableName, valueField.Interface(), condition, t)
			default:
				// apply value
				t.Value = valueField.Interface()

				// apply json name as column name if not specified
				if t.Column == "" {
					column := utils.StructFieldJSONTag(typeField)
					t.Column = column
				}
				ParseSearchTag(tableName, valueField.Interface(), condition, parent...)
			}
		}
	}
}

// makeTag 解析search的tag标签
func makeTag(tableName string, tag string) *Tag {
	r := &Tag{
		Table: tableName,
	}
	tags := strings.Split(tag, ";")
	var ts []string
	for _, t := range tags {
		ts = strings.Split(t, ":")
		if len(ts) == 0 {
			continue
		}
		switch ts[0] {
		case "op":
			if len(ts) > 1 {
				r.Operator = ts[1]
			}
		case "col", "column":
			if len(ts) > 1 {
				r.Column = ts[1]
			}
		case "table":
			if len(ts) > 1 {
				r.Table = ts[1]
			}
		case "on": // e.g. "on:a=b,c<>d"
			if len(ts) > 1 {
				cond := strings.Split(ts[1], ",")
				on := make([][]string, len(cond))
				for i, v := range cond {
					comp := strings.Split(v, "=")
					if len(comp) == 2 {
						on[i] = []string{comp[0], "=", comp[1]}
					}
					comp = strings.Split(v, "<>")
					if len(comp) == 2 {
						on[i] = []string{comp[0], "<>", comp[1]}
					}
					comp = strings.Split(v, ">")
					if len(comp) == 2 {
						on[i] = []string{comp[0], ">", comp[1]}
					}
					comp = strings.Split(v, ">=")
					if len(comp) == 2 {
						on[i] = []string{comp[0], ">=", comp[1]}
					}
					comp = strings.Split(v, "<=")
					if len(comp) == 2 {
						on[i] = []string{comp[0], "<=", comp[1]}
					}
				}

				if len(on) > 0 {
					r.On = on
				}
			}
		case "onraw":
			if len(ts) > 1 {
				r.OnRaw = ts[1]
			}
		case "join":
			if len(ts) > 1 {
				r.Join = ts[1]
			}
		case "as":
			if len(ts) > 1 {
				r.As = ts[1]
			}
		case "jsonvalue":
			if len(ts) > 1 {
				r.JsonValue = ts[1]
			}
		case "skipcheckempty":
			r.SkipCheckEmpty = true
		}
	}
	return r
}
