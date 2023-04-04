package validate

import (
	"fmt"
	"regexp"
)

func SAPVirtualInstanceName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[A-Z][A-Z0-9][A-Z0-9]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must be an alphabet for the first letter. The second and third letters must be alphanumeric and all alphabets must be uppercase", k))
		return warnings, errors
	}

	return warnings, errors
}
