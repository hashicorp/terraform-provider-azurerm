package validate

import "fmt"

func ValidateEventHubPartitionCount(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)

	if !(1024 >= value && value >= 1) {
		errors = append(errors, fmt.Errorf("EventHub Partition Count has to be between 1 and 32 or between 1 and 1024 if using a dedicated Event Hubs Cluster"))
	}

	return warnings, errors
}
