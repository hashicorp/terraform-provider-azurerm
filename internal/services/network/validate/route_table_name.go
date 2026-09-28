// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RouteTableName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,78}[a-zA-Z0-9_]?$`), "should be between 1 and 80 characters, start with an alphanumeric, end with an alphanumeric or underscore and can contain alphanumerics, underscores, periods, and hyphens")(v, k)
}
