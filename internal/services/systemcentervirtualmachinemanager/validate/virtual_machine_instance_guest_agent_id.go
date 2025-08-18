// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/parse"
)

func SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := parse.SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
