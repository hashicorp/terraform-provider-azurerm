// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SpringCloudCustomDomainName(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([a-z0-9]+(-[a-z0-9]+)*\.)+[a-z]{2,}$`), "should be a valid domain name")(input, key)
}
