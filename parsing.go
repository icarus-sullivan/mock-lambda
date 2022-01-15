package lambda

import (
	"encoding/json"
	"strings"
)

func sanitize(src string) string {
	// When we get this env value it's escaped in a way javascript understands but not
	// go
	// 	- unescape to a single backslash
	return strings.ReplaceAll(src, "\\\\", "\\")
}

func decode(src string, target interface{}) error {
	sanitizedEvent := sanitize(src)
	err := json.Unmarshal([]byte(sanitizedEvent), target)
	if err != nil {
		return err
	}
	return nil
}
