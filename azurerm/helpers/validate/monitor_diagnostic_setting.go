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

	if len(value) < 1 || len(value) > 260 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 260 characters: %q", k, value))
	}

	return warnings, errors
}
