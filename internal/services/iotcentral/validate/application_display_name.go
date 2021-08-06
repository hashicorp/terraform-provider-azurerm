package validate

import (
	"fmt"
	"regexp"
)

func ApplicationDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^.{1,200}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("test: %s, %q length should between 1~200", k, v))
	}
	return warnings, errors
}
