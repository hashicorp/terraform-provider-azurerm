// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FluidRelayServerName(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[-0-9a-zA-Z]{1,50}$`), "should contain only alphanumeric characters and hyphens, up to 50 characters long")(input, key)
}
