package validate

import (
	"fmt"
	"regexp"
)

func StorageAccountName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile(`\A([a-z0-9]{3,24})\z`).MatchString(input) {
		errors = append(errors, fmt.Errorf("name (%q) can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long", input))
	}

	return warnings, errors
}
