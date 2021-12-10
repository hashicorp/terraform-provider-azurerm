package validate

import (
	"fmt"
	"strings"
)

func FileNameFormat(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	requiredComponents := []string{
		"{iothub}",
		"{partition}",
		"{YYYY}",
		"{MM}",
		"{DD}",
		"{HH}",
		"{mm}",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(value, component) {
			errors = append(errors, fmt.Errorf("%s needs to contain %q", k, component))
		}
	}

	return warnings, errors
}
