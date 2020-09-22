package azure

import (
	"fmt"
	"regexp"
)

// The alert name can contain any characters, but the following characters: <, >, *, %, &, :, \, ?, +, / and can't have more than 260 characters.
func ValidateDiagnosticSettingsName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if matched := regexp.MustCompile(`^[^<>*%&:\\\/?]{0,260}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf(`%q can't contain '<,>,*,%%,&,:,\,/,?' or control characters, and can't have more than 260 characters.`, k))
	}

	return warnings, errors
}
