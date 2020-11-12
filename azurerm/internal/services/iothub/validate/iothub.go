package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/parse"
)

func IotHubID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.IotHubID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func IoTHubName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile(`^[0-9a-zA-Z-]{1,}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	return warnings, errors
}

func IoTHubConsumerGroupName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: - . _
	if matched := regexp.MustCompile(`^[0-9a-zA-Z-._]{1,}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes, periods and underscores", k))
	}

	return warnings, errors
}

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

func IotHubSharedAccessPolicyName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The name attribute rules are :
	// 1. must not be empty.
	// 2. must not exceed 64 characters in length.
	// 3. can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens

	if !regexp.MustCompile(`[a-zA-Z0-9!._-]{1,64}`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must not be empty, and must not exceed 64 characters in length, and can only contain alphanumeric characters, exclamation marks, periods, underscores and hyphens", k))
	}

	return nil, errors
}
