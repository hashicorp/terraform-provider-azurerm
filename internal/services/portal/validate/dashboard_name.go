package validate

import (
	"fmt"
	"regexp"
)

func DashboardName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q may not exceed 64 characters in length", k))
	}

	// only alpanumeric and hyphens
	if matched := regexp.MustCompile(`^[-\w]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric and hyphen characters", k))
	}

	return warnings, errors
}
