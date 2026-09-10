// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ElasticSanName(i interface{}, k string) ([]string, []error) {
	return elasticSanResourceName(24)(i, k)
}

// all elastic san resource names must be 3 to maxLength characters
func elasticSanResourceName(maxLength int) func(interface{}, string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(
			regexp.MustCompile(fmt.Sprintf(`^[a-z0-9][a-z0-9_-]{1,%d}[a-z0-9]$`, maxLength-2)),
			fmt.Sprintf("must be between 3 and %d characters. It can contain only lowercase letters, numbers, underscores (_) and hyphens (-). It must start and end with a lowercase letter or number", maxLength),
		),
		validation.StringDoesNotMatch(regexp.MustCompile(`[_-][_-]`), "must have hyphens and underscores be surrounded by alphanumeric character"),
	)
}
