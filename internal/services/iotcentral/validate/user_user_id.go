// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func UserUserID(i interface{}, k string) (warnings []string, errors []error) {
	id, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return warnings, errors
	}

	err := validateUserId(id)
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func validateUserId(id string) error {
	formatPattern := `^[a-zA-Z0-9-_]*$` // https://learn.microsoft.com/en-us/rest/api/iotcentral/dataplane/users/get?view=rest-iotcentral-dataplane-2022-07-31&tabs=HTTP
	formatRegex, err := regexp.Compile(formatPattern)
	if err != nil {
		return fmt.Errorf("error compiling format regex: %s error: %+v", formatPattern, err)
	}

	if !formatRegex.MatchString(id) {
		return fmt.Errorf("iot central userId %q is invalid, regex pattern: ^[a-zA-Z0-9-_]*$", id)
	}

	// Ensure the string does not start with a hyphen.
	// Solves for (?!-)
	startHyphenPattern := `^-`
	startHyphenRegex, err := regexp.Compile(startHyphenPattern)
	if err != nil {
		return fmt.Errorf("error compiling start hyphen regex: %s error: %+v", startHyphenPattern, err)
	}

	if startHyphenRegex.MatchString(id) {
		return fmt.Errorf("iot central userId %q is invalid, regex pattern: ^[a-zA-Z0-9-_]*$", id)
	}

	return nil
}
