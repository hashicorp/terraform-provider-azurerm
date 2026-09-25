// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TimeInterval(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([0-9][0-9]):([0-5][0-9]):([0-5][0-9])$`), "must be in the form HH:MM:SS between 00:00:00 and 99:59:59")(i, k)
}
