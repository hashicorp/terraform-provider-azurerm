package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func NextGenerationFirewallName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRuleStackName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func paloAltoNameValidation(input interface{}, k string) (warnings []string, errors []error) {
	value := input.(string)

	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{1,128}$`).Match([]byte(value)); !matched {
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes, and must be between 1 and 128 characters in length", k))
	}

	if strings.HasSuffix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("%q cannot start or end with a `-`"))
	}

	return
}
