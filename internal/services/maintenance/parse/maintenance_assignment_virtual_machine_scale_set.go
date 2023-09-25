// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type MaintenanceAssignmentVirtualMachineScaleSetId struct {
	VirtualMachineScaleSetId    *commonids.VirtualMachineScaleSetId
	VirtualMachineScaleSetIdRaw string
	Name                        string
}

func MaintenanceAssignmentVirtualMachineScaleSetID(input string) (*MaintenanceAssignmentVirtualMachineScaleSetId, error) {
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine Scale Set ID (%q)", input)
	}

	targetResourceId, name := groups[1], groups[2]
	virtualMachineScaleSetId, err := commonids.ParseVirtualMachineScaleSetIDInsensitively(targetResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine Scale Set ID: %q: Expected valid virtual machine scale set ID", input)
	}

	return &MaintenanceAssignmentVirtualMachineScaleSetId{
		VirtualMachineScaleSetId:    virtualMachineScaleSetId,
		VirtualMachineScaleSetIdRaw: targetResourceId,
		Name:                        name,
	}, nil
}
