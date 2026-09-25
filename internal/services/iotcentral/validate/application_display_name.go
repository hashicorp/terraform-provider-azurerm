// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationDisplayName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^.{1,200}$`), "length should between 1~200")(v, k)
}
