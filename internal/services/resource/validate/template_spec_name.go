// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TemplateSpecName(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\w\.\-\(\)]{1,64}$`), "must only contain alpha-numeric characters, parenthesis, underscores, dashes and periods and be between 1 and 64 characters in length")(input, key)
}
