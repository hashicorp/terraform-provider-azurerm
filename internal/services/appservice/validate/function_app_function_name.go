// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FunctionAppFunctionName(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-zA-Z](([-_0-9a-zA-Z-]{0,126})[-_0-9a-zA-Z])?$`), "must start with a letter, may only contain alphanumeric characters, dashes, underscore and up to 128 characters in length")(input, key)
}
