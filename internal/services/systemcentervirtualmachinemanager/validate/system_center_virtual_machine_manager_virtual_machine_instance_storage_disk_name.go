// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDiskName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[a-zA-Z]([-_a-zA-Z0-9]{0,78}[a-zA-Z0-9])?$"), "must start and end with a letter, may contain alphanumeric characters, dashes or underscores and must be between 1 and 80 characters long")(i, k)
}
