// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationName(v interface{}, k string) ([]string, []error) {
	// Portal: The value must contain only alphanumeric characters or the following: -
	return validation.StringMatch(regexp.MustCompile(`^[a-z\d][a-z\d-]{0,61}[a-z\d]$`), "may only contain lowercase alphanumeric characters and dashes, length between 2-63, and must start and end with an alphanumeric character")(v, k)
}
