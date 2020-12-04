package validate

import (
	"fmt"
	"regexp"
)

func TemplateDeploymentContentVersion(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	var errors []error
	if matched := regexp.MustCompile(`^[0-9]+.[0-9]+.[0-9]+.[0-9]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q only contains numbers and dots", k))
	}

	return nil, errors
}
