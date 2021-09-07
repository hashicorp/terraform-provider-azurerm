package azure

import (
	"fmt"
	"strings"
)

// QuotedStringSlice formats a string slice into a quoted string containing all segments passed in a slice (e.g. string[]{"one", "two", "three"} will return {"one", "two" or "three"}). Useful for error messages with multiple possible values.
func QuotedStringSlice(strs []string) string {
	if len(strs) == 1 {
		return fmt.Sprint(`"`, strs[0], `"`)
	}

	var sb strings.Builder

	for i, str := range strs {
		if i < (len(strs) - 1) {
			if i == (len(strs) - 2) {
				sb.WriteString(fmt.Sprint(`"`, str, `"`))
			} else {
				sb.WriteString(fmt.Sprint(`"`, str, `", `))
			}
		} else {
			sb.WriteString(fmt.Sprint(` or "`, str, `"`))
		}
	}

	return sb.String()
}
