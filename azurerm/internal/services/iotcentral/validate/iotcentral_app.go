package validate

import (
	"fmt"
	"regexp"
)

func IotCentralAppName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile(`^[a-z\d][a-z\d-]{0,61}[a-z\d]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}
	return warnings, errors
}

func IotCentralAppSubdomain(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// Portal: The value must contain only alphanumeric characters or the following: -
	if matched := regexp.MustCompile(`^[a-z\d][a-z\d-]{0,61}[a-z\d]$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}
	return warnings, errors
}

func IotCentralAppDisplayName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^.{1,200}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q length should between 1~200", k))
	}
	return warnings, errors
}

func IotCentralAppTemplateName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^.{1,50}$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q length should between 1~50", k))
	}
	return warnings, errors
}
