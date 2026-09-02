// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ElasticSanVolumeName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{1,61}[a-z0-9]$`), "must be between 3 and 63 characters. It can contain only lowercase letters, numbers, underscores (_) and hyphens (-). It must start and end with a lowercase letter or number"),
		validation.StringDoesNotMatch(regexp.MustCompile(`[_-][_-]`), "must have hyphens and underscores be surrounded by alphanumeric character"),
	)(i, k)
}
