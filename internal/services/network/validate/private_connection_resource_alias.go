// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func PrivateConnectionResourceAlias(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(
		regexp.MustCompile(`\.azure\.privatelinkservice$`),
		"expected to have suffix `.azure.privatelinkservice`",
	)(input, key)
}
