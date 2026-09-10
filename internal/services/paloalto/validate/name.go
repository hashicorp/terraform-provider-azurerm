// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func NextGenerationFirewallName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRuleStackName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRuleStackCertificateName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRuleStackFQDNListName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func LocalRuleStackRuleName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func DestinationNATName(input interface{}, k string) (warnings []string, errors []error) {
	return paloAltoNameValidation(input, k)
}

func paloAltoNameValidation(input interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-]{1,128}$`), "may only contain alphanumeric characters and dashes, and must be between 1 and 128 characters in length"),
		validation.StringDoesNotMatch(regexp.MustCompile(`^-|-$`), "cannot start or end with a `-`"),
	)(input, k)
}
