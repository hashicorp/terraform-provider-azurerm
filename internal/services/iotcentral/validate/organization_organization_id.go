// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

func OrganizationOrganizationID(i interface{}, k string) (warnings []string, errors []error) {
	id, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to be a string", k))
		return warnings, errors
	}

	err := validateOrganizationId(id)
	if err != nil {
		errors = append(errors, err)
		return warnings, errors
	}

	return warnings, errors
}

func validateOrganizationId(id string) error {
	// Ensure the string follows the desired format.
	// Regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$
	// The negative lookahead (?!-) is not supported in Go's standard regexp package
	formatPattern := `^[a-z0-9-]{1,48}[a-z0-9]$`
	formatRegex, err := regexp.Compile(formatPattern)
	if err != nil {
		return fmt.Errorf("error compiling format regex: %s error: %+v", formatPattern, err)
	}

	if !formatRegex.MatchString(id) {
		return fmt.Errorf("iot central organizationId %q is invalid, regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$", id)
	}

	// Ensure the string does not start with a hyphen.
	// Solves for (?!-)
	startHyphenPattern := `^-`
	startHyphenRegex, err := regexp.Compile(startHyphenPattern)
	if err != nil {
		return fmt.Errorf("error compiling start hyphen regex: %s error: %+v", startHyphenPattern, err)
	}

	if startHyphenRegex.MatchString(id) {
		return fmt.Errorf("iot central organizationId %q is invalid, regex pattern: ^(?!-)[a-z0-9-]{1,48}[a-z0-9]$", id)
	}

	return nil
}
