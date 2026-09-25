// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BotChannelRegistrationIconUrl(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`\.png$`), "only png is supported")(i, k)
}
