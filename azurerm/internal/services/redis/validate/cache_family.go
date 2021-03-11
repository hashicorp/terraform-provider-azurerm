package validate

import (
	"fmt"
	"strings"
)

func CacheFamily(v interface{}, _ string) (warnings []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"c": true,
		"p": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Family can only be C or P"))
	}
	return warnings, errors
}
