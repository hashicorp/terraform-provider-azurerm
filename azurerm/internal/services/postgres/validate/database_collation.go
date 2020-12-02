package validate

import (
	"fmt"
	"regexp"
)

func DatabaseCollation(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	matched, _ := regexp.MatchString(`^[-A-Za-z0-9_. ]+$`, v)

	if !matched {
		errors = append(errors, fmt.Errorf("%s contains invalid characters, only alphanumeric, underscore, space or hyphen characters are supported, got %s", k, v))
		return
	}

	return warnings, errors
}
