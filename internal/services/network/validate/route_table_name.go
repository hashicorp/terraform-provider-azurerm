package validate

import (
	"fmt"
	"regexp"
)

func RouteTableName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_]?$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q should be between 1 and 80 characters, start with an alphanumeric, end with an alphanumeric or underscore and can contain alphanumerics, underscores, periods, and hyphens", k))
	}

	return warnings, errors
}
