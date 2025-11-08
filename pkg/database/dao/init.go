package dao

import "fmt"

var enclose = "`"

// SetEnclose postgres enclocs is ", other is `
func SetEnclose(e string) {
	enclose = e
}

// GetEnclose get enclose
func GetEnclose() string {
	return enclose
}

// Enclose enclose string
func Enclose(s string) string {
	return fmt.Sprintf("%s%s%s", enclose, s, enclose)
}

// EncloseTableColumn enclose table column
func EncloseTableColumn(t string, s string) string {
	return fmt.Sprintf("%s%s%s.%s%s%s", enclose, t, enclose, enclose, s, enclose)
}

func EncloseTableJson(t, c string) string {
	return fmt.Sprintf(`%s.%s`, t, c)
}
