package validate

import (
	"fmt"
	"regexp"
)

func DataboxEdgePhoneNumber(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^([\+])?(\d{1,2})?[\s-]?\(?(\d{3})\)?[\s-]?(\d{3})[-](\d{4})$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q may contain parentheses, hyphens, plus sign or digits only, must be in a valid 10 digit phone number format with country code being optional(e.g. 123 555-6789)", k))
	}

	return warnings, errors
}
