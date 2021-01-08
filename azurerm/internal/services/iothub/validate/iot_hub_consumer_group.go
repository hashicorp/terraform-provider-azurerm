package validate

import (
	"fmt"
	"regexp"
)

func IoTHubConsumerGroupName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: - . _
	if matched := regexp.MustCompile(`^[0-9a-zA-Z-._]{1,}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes, periods and underscores", k))
	}

	return warnings, errors
}
