// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

// DataCollectionRuleAssociationName validates that the name does not contain control characters or specific forbidden characters
func DataCollectionRuleAssociationName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	// Check for control characters (ASCII 0-31 and 127)
	controlCharPattern := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	if controlCharPattern.MatchString(v) {
		errors = append(errors, fmt.Errorf("property `%s` cannot contain control characters", k))
		return warnings, errors
	}

	// Check for forbidden characters: < > % & : \ ? /
	// Use regex to find all forbidden characters at once
	forbiddenCharPattern := regexp.MustCompile(`[<>%&:\\?/]`)
	if forbiddenCharPattern.MatchString(v) {
		errors = append(errors, fmt.Errorf("property `%s` cannot contain any of the following characters: `<`, `>`, `%%`, `&`, `:`, `\\`, `?`, `/`", k))
		return warnings, errors
	}

	return warnings, errors
}
