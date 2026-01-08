// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

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
			errors = append(errors, fmt.Errorf("the reserved endpoint name %s could not be used as a name for a custom endpoint", name))
		}
	}

	return warnings, errors
}
