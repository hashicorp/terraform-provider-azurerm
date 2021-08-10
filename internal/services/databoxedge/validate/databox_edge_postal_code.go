package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgePostalCode(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^([\d]{5})((-)([\d]{4}))?$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must consist of 5 digits unless it is in the ZIP+4 format then it must consist of five digits, a hyphen, then four digits", k))
	}

	return warnings, errors
}
