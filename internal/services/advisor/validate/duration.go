// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func Duration(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^(?:[0-9]{1,2}:)?[0-9]{2}:[0-9]{2}:[0-9]{2}$`), "must be in format DD:HH:MM:SS. If DD is 00, it has to be omit")(v, k)
}
