package lambda

import "strings"

func SanitizeJSON(json string) string {
	// When we get this env value it's escaped in a way javascript understands but not
	// go
	// 	- unescape to a single backslash
	return strings.ReplaceAll(json, "\\\\", "\\")
}
