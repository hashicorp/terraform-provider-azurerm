package validate

import (
	"fmt"
	"strings"
)

func prettyErrorString(strs []string) string {
	if len(strs) == 1 {
		return fmt.Sprint("\"", strs[0], "\"")
	}

	var sb strings.Builder

	for i, str := range strs {
		if i < (len(strs) - 1) {
			if i == (len(strs) - 2) {
				sb.WriteString(fmt.Sprint("\"", str, "\""))
			} else {
				sb.WriteString(fmt.Sprint("\"", str, "\", "))
			}
		} else {
			sb.WriteString(fmt.Sprint(" or \"", str, "\""))
		}
	}

	return sb.String()
}
