// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationTemplateName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^.{1,50}$`), "length should between 1~50")(v, k)
}
