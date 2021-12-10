package validate

import (
	"fmt"
	"strings"
)

func ValidateEventHubArchiveNameFormat(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	requiredComponents := []string{
		"{Namespace}",
		"{EventHub}",
		"{PartitionId}",
		"{Year}",
		"{Month}",
		"{Day}",
		"{Hour}",
		"{Minute}",
		"{Second}",
	}

	for _, component := range requiredComponents {
		if !strings.Contains(value, component) {
			errors = append(errors, fmt.Errorf("%s needs to contain %q", k, component))
		}
	}

	return warnings, errors
}
