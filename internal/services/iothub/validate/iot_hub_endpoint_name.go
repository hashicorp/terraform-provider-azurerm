package validate

import "fmt"

func IoTHubEndpointName(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(string)

	reservedNames := []string{
		"events",
		"operationsMonitoringEvents",
		"fileNotifications",
		"$default",
	}

	for _, name := range reservedNames {
		if name == value {
			errors = append(errors, fmt.Errorf("The reserved endpoint name %s could not be used as a name for a custom endpoint", name))
		}
	}

	return warnings, errors
}
