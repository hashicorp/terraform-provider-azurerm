// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DomainServiceName(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^(([0-9a-zA-Z])|(([0-9a-zA-Z][0-9a-zA-Z-]{0,28}[0-9a-zA-Z])))(\.[0-9a-zA-Z-]+)+$`), "must be a valid FQDN and the first element must be 15 or fewer characters")(input, key)
}
