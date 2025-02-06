// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

type MaintenanceAssignmentVirtualMachineId struct {
	VirtualMachineId    *commonids.VirtualMachineId
	VirtualMachineIdRaw string
	Name                string
}

func MaintenanceAssignmentVirtualMachineID(input string) (*MaintenanceAssignmentVirtualMachineId, error) {
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine ID (%q)", input)
	}

	targetResourceId, name := groups[1], groups[2]
	virtualMachineId, err := commonids.ParseVirtualMachineIDInsensitively(targetResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine ID: %q: Expected valid virtual machine ID", input)
	}

	return &MaintenanceAssignmentVirtualMachineId{
		VirtualMachineId:    virtualMachineId,
		VirtualMachineIdRaw: targetResourceId,
		Name:                name,
	}, nil
}
