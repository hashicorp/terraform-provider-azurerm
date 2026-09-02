// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ElasticSanName(i interface{}, k string) ([]string, []error) {
	return elasticSanResourceName(3, 24)(i, k)
}

func elasticSanResourceName(minLength, maxLength int) func(interface{}, string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(
			regexp.MustCompile(fmt.Sprintf(`^[a-z0-9][a-z0-9_-]{%d,%d}[a-z0-9]$`, minLength-2, maxLength-2)),
			fmt.Sprintf("must be between %d and %d characters. It can contain only lowercase letters, numbers, underscores (_) and hyphens (-). It must start and end with a lowercase letter or number", minLength, maxLength),
		),
		validation.StringDoesNotMatch(regexp.MustCompile(`[_-][_-]`), "must have hyphens and underscores be surrounded by alphanumeric character"),
	)
}
