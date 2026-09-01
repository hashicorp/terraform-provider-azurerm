// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BatchMaxWaitTime(input interface{}, key string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`[0-9]{2}:[0-9]{2}:[0-9]{2}`), "must have the following format hh:mm:ss")(input, key)
}
