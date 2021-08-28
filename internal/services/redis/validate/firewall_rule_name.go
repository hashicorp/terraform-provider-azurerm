package validate

import (
	"fmt"
	"regexp"
)

func FirewallRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^\w+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and underscores", k))
	}

	return warnings, errors
}
