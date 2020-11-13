package validate

import (
	"fmt"
	"regexp"
)

func DnsZoneSOARecordEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[a-zA-Z\d_-][a-zA-Z\d._-]{0,62}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must be between 1 and 63 characters in length and contains only letters, numbers, underscores, dashes and periods", k))
		return warnings, errors
	}

	return warnings, errors
}
