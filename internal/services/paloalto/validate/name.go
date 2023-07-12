package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func NextGenerationFirewallName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRulestackName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRulestackCertificateName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRulestackFQDNListName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRulestackRuleName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func paloAltoNameValidation(input interface{}, k string) (warnings []string, errors []error) {
	value, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be of type string", k))
		return
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]{1,128}$`).Match([]byte(value)); !matched {
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes, and must be between 1 and 128 characters in length", k))
	}

	if strings.HasSuffix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("%q cannot start or end with a `-`"))
	}

	return
}
