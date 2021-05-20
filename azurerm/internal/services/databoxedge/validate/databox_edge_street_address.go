package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgeStreetAddress(v []interface{}, k string) (warnings []string, errors []error) {
	if len(v) == 0 {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return warnings, errors
	}

	for i, addressLine := range v {
		if !regexp.MustCompile(`^[\s\S]{1,35}$`).MatchString(addressLine.(string)) {
			errMsg := fmt.Sprintf("'shipping_info' %q line %d must be between 1 and 35 characters in length", k, (i + 1))

			if len(errors) > 0 {
				errors = append(errors, fmt.Errorf("\n        %s", errMsg))
			} else {
				errors = append(errors, fmt.Errorf("%s", errMsg))
			}
		}
	}

	return warnings, errors
}
