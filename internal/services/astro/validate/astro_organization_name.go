package validate

import (
	"fmt"
	"regexp"
)

func OrganizationAndOrganizationWorkspaceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 50 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 50 characters in length", k))
	}

	regex := `^[a-zA-Z0-9][a-zA-Z0-9_\-.: ]*$`
	if matched := regexp.MustCompile(regex).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must match the regular expression `%s`", k, regex))
	}
	return warnings, errors
}
