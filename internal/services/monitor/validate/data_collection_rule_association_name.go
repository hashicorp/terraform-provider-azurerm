// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
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
	forbiddenChars := "<>%&:\\?/"
	for _, char := range forbiddenChars {
		if strings.ContainsRune(v, char) {
			errors = append(errors, fmt.Errorf("property `%s` cannot contain the character `%c`", k, char))
			return warnings, errors
		}
	}

	return warnings, errors
}
