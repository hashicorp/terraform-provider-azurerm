package validate

import (
	"fmt"
	"regexp"
)

func MonitorDiagnosticSettingName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if regexp.MustCompile(`[<>*%&:\\?+\/]+`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"characters <, >, *, %%, &, :, \\, ?, +, / are not allowed in %q: %q", k, value))
	}

	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must be longer than 0 characters: %q %d", k, value, len(value)))
	}

	if len(value) > 260 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 260 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}
