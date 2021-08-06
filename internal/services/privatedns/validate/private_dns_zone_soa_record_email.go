package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func PrivateDnsZoneSOARecordEmail(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		return warnings, errors
	}

	vSegments := strings.Split(value, ".")
	if len(vSegments) < 2 || len(vSegments) > 34 {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 34 segments", k))
		return warnings, errors
	}

	for _, segment := range vSegments {
		if segment == "" {
			errors = append(errors, fmt.Errorf("%q cannot contain consecutive period", k))
			return warnings, errors
		}

		if len(segment) > 63 {
			errors = append(errors, fmt.Errorf("the each segment of the `email` must contain between 1 and 63 characters"))
			return warnings, errors
		}
	}

	if !regexp.MustCompile(`^[a-zA-Z\d._-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q only contains letters, numbers, underscores, dashes and periods", k))
		return warnings, errors
	}

	return warnings, errors
}
