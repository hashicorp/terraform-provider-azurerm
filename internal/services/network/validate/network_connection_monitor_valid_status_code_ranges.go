package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func NetworkConnectionMonitorValidStatusCodeRanges(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, value))
		return warnings, errors
	}

	if len(value) != 7 && len(value) != 3 {
		errors = append(errors, fmt.Errorf("The len of %q should be 3 or 7: %q", k, value))
		return warnings, errors
	}

	// Here the format of the expected code range is `301-304`
	if len(value) == 7 {
		if !regexp.MustCompile(`^([1-5][0-9][0-9]-([1-5][0-9][0-9]|600))$`).MatchString(value) {
			errors = append(errors, fmt.Errorf("%q can contain hyphen: %q", k, value))
			return warnings, errors
		} else {
			vArray := strings.Split(value, "-")

			startNumber, err := strconv.Atoi(vArray[0])
			if err != nil {
				errors = append(errors, fmt.Errorf("expected %s on the left of - to be an integer, got %v: %v", k, value, err))
				return warnings, errors
			}

			endNumber, err := strconv.Atoi(vArray[1])
			if err != nil {
				errors = append(errors, fmt.Errorf("expected %s on the right of - to be an integer, got %v: %v", k, value, err))
				return warnings, errors
			}

			if startNumber >= endNumber {
				errors = append(errors, fmt.Errorf("the start number of %q should less than the end number: %q", k, value))
				return warnings, errors
			}
		}
	}

	// Here the format of the expected code ranges are `2xx` and `418`
	if len(value) == 3 {
		if !regexp.MustCompile(`^([1-5][0-9x][0-9x]|600)$`).MatchString(value) {
			errors = append(errors, fmt.Errorf("%q can contain number with x or pure number: %q", k, value))
			return warnings, errors
		}
	}

	return warnings, errors
}
