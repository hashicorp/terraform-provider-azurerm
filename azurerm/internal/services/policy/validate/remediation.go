package validate

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
)

func RemediationName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The service returns error when name of remediation is too long
	// error: The remediation name cannot be empty and must not exceed '260' characters.
	// By my additional test, the name of remediation cannot contain the following characters: %^#/\&?.
	if len(v) == 0 || len(v) > 260 {
		errors = append(errors, fmt.Errorf("%s cannot be empty and must not exceed '260' characters", k))
		return
	}
	const invalidCharacters = `%^#/\&?`
	if strings.ContainsAny(v, invalidCharacters) {
		errors = append(errors, fmt.Errorf("%s cannot contain the following characters: %s", k, invalidCharacters))
	}
	// Despite the service accepts remediation name with capitalized characters, but in the response,
	// all upper case characters will be converted to lower cases. Therefore we forbid user to use upper case letters here
	if v != strings.ToLower(v) {
		errors = append(errors, fmt.Errorf("%s cannot contain upper case letters", k))
	}

	return warnings, errors
}

func RemediationID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.PolicyRemediationID(v); err != nil {
		errors = append(errors, fmt.Errorf("cannot parse %q as a Policy Remediation ID: %+v", k, err))
		return
	}

	return warnings, errors
}
