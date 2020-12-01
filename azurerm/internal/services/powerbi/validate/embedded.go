package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-uuid"
)

func PowerBIEmbeddedName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-z][a-z0-9]{3,63}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 4 and 64 characters in length and contains only lowercase letters or numbers.", k))
	}

	return warnings, errors
}

func PowerBIEmbeddedAdministratorName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`).MatchString(value) {
		if _, err := uuid.ParseUUID(value); err != nil {
			errors = append(errors, fmt.Errorf("%q isn't a valid email address.", k))
			errors = append(errors, fmt.Errorf("%q isn't a valid UUID (%q): %+v", k, v, err))
		}
	}

	return warnings, errors
}
