package validate

import (
	"fmt"
	"regexp"
)

// JobName validates the name of a Batch job
func JobName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w-]{1,64}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q:%q can only contain any combination of alphanumeric characters along with dash (-) and underscore (_). The name must be from 1 through 64 characters long.", k, value))
	}

	return warnings, errors
}
