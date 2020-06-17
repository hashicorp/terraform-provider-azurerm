package parse

import (
	"fmt"
	"regexp"

	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

type MaintenanceAssignmentVirtualMachineId struct {
	VirtualMachineId    *parseCompute.VirtualMachineId
	VirtualMachineIdRaw string
	Name                string
}

func MaintenanceAssignmentVirtualMachineID(input string) (*MaintenanceAssignmentVirtualMachineId, error) {
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine ID (%q)", input)
	}

	targetResourceId, name := groups[1], groups[2]
	virtualMachineId, err := parseCompute.VirtualMachineID(targetResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine ID: %q: Expected valid virtual machine ID", input)
	}

	return &MaintenanceAssignmentVirtualMachineId{
		VirtualMachineId:    virtualMachineId,
		VirtualMachineIdRaw: targetResourceId,
		Name:                name,
	}, nil
}
