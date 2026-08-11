// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func HDInsightClusterVersion(i interface{}, k string) (warnings []string, errors []error) {
	// 3.6, 3333.6666 or 1.2.3000.45
	return validation.Any(
		// `major minor`
		validation.StringMatch(regexp.MustCompile(`^(\d)+(.){1}(\d)+$`), "must be a version in the format `x.y` or `a.b.c.d`"),
		// `major minor build release`
		validation.StringMatch(regexp.MustCompile(`^(\d)+(.)(\d)+(.)(\d)+(.)(\d)+$`), "must be a version in the format `x.y` or `a.b.c.d`"),
	)(i, k)
}

func HDInsightName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	// The name must be 59 characters or less and can contain letters, numbers, and hyphens (but the first and last character must be a letter or number).
	if matched := regexp.MustCompile(`(^[a-zA-Z0-9])([a-zA-Z0-9-]{1,57})([a-zA-Z0-9]$)`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q must be 59 characters or less and can contain letters, numbers, and hyphens (but the first and last character must be a letter or number)", k))
	}

	return warnings, errors
}
