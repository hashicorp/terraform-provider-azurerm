// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TriggerTimespan(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^\-?((\d+)\.)?(\d\d):(60|([0-5][0-9])):(60|([0-5][0-9]))`), "invalid timespan, must be of format hh:mm:ss")(i, k)
}
