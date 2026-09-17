// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TimeOfDayInUTC(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]|[0-9]):([0-5][0-9])$"), "must match the format HHmm where HH is 00-23 and mm is 00-59")(i, k)
}
