package validate

import (
	"fmt"
	"regexp"
)

func ElasticSanName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't!", k))
		return
	}

	// name can be 3-24 characters in length
	const minLength = 3
	if len(v) < minLength {
		errors = append(errors, fmt.Errorf("%q can be at least %d characters, got %d", k, minLength, len(v)))
	}

	const maxLength = 24
	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	if matched := regexp.MustCompile(`^[a-z0-9_-]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain lower case characters, numbers, dashes and underscores", k))
	}

	if matched := regexp.MustCompile(`^[a-z0-9]`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must begin with an alphanumeric character", k))
	}

	if matched := regexp.MustCompile(`[a-z0-9]$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must end with an alphanumeric character", k))
	}

	if matched := regexp.MustCompile(`(_|-)+[_-]`).Match([]byte(v)); matched {
		errors = append(errors, fmt.Errorf("%q must have hyphens and underscores be surrounded by letters or numbers", k))
	}

	return warnings, errors
}
