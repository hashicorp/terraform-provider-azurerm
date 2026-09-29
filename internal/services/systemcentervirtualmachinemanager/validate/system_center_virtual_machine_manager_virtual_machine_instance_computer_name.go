// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SystemCenterVirtualMachineManagerVirtualMachineInstanceComputerName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9]{1,}$"), "must only contain alphanumeric characters")(i, k)
}
