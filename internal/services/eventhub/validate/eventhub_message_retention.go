package validate

import "fmt"

func ValidateEventHubMessageRetentionCount(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)

	if !(90 >= value && value >= 1) {
		errors = append(errors, fmt.Errorf("EventHub Retention Count has to be between 1 and 7 or between 1 and 90 if using a dedicated Event Hubs Cluster"))
	}

	return warnings, errors
}
