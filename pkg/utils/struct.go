package utils

import (
	"reflect"
	"strings"
)

// StructFieldJSONTag StructFieldJSONTag
func StructFieldJSONTag(t reflect.StructField) string {
	switch jsonTag := t.Tag.Get("json"); jsonTag {
	case "-":
	case "":
		return ""
	default:
		parts := strings.Split(jsonTag, ",")
		name := parts[0]
		if name == "" {
			name = t.Name
		}
		return name
	}

	return t.Name
}

// IsStruct .
func IsStruct(v reflect.Value) bool {
	if v.Kind() == reflect.Pointer {
		return IsStruct(v.Elem())
	}

	return v.Kind() == reflect.Struct
}
