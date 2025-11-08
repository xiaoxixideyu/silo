package utils

import (
	"encoding/json"
	"errors"
	"math"
	"net/url"
	"reflect"
)

// MapKeys get keys from map
func MapKeys[T any](v map[string]T) []string {
	rv := reflect.ValueOf(v)
	var result []string
	for _, kv := range rv.MapKeys() {
		result = append(result, kv.String())
	}
	return result
}

// MapValues MapValues
func MapValues[T any](m map[string]T) []T {
	r := []T{}
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

// ArrayElementExists check wether element exists in the array
func ArrayElementExists(source any, element any) bool {
	arr := reflect.ValueOf(source)

	if arr.Kind() != reflect.Array && arr.Kind() != reflect.Slice {
		panic("ArrayElementExists invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == element {
			return true
		}
	}

	return false
}

// ToJSON ToJSON
func ToJSON(b any) (string, error) {
	res, err := json.Marshal(b)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// AnyToArray convert string or []string to []string
func AnyToArray[T any](i any) []T {
	typeOf := reflect.TypeOf(i)
	k := typeOf.Kind()
	switch k {
	case reflect.Array, reflect.Slice:
		r := []T{}
		source := reflect.ValueOf(i).Interface()
		for _, v := range source.([]any) {
			r = append(r, v.(T))
		}
		return r
	default:
		return []T{reflect.ValueOf(i).Interface().(T)}
	}
}

// IsNil IsNil
func IsNil(input any) bool {
	if input == nil {
		return true
	}
	kind := reflect.ValueOf(input).Kind()
	switch kind {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan:
		return reflect.ValueOf(input).IsNil()
	default:
		return false
	}
}

// SafeInt64ToInt safe int64 to int
func SafeInt64ToInt(value int64) (int, error) {
	if value > math.MaxInt || value < math.MinInt {
		return 0, errors.New("int64 value out of int range")
	}
	return int(value), nil
}

// append slice with distinct elements
func AppendDistinct[T comparable](src []T, dst T) []T {
	if !Contains[T](src, dst) {
		src = append(src, dst)
	}
	return src
}

// JoinURL join url
func JoinURL(baseURL string, path string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	return u.JoinPath(path).String(), nil
}
