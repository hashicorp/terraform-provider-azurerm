// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SystemCenterVirtualMachineManagerVirtualMachineInstanceMacAddress(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$"), "must be in format `00:00:00:00:00:00`")(i, k)
}
