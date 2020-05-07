package azure

import (
	"strings"
)

func StringContains(sourceStr, subStr string) bool {
	sourceStr, subStr = strings.ToUpper(sourceStr), strings.ToUpper(subStr)
	return strings.Contains(sourceStr, subStr)
}
