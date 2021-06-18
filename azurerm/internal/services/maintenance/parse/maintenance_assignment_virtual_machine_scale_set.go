package parse

import (
	"fmt"
	"regexp"

	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
)

type MaintenanceAssignmentVirtualMachineScaleSetId struct {
	VirtualMachineScaleSetId    *parseCompute.VirtualMachineScaleSetId
	VirtualMachineScaleSetIdRaw string
	Name                        string
}

func MaintenanceAssignmentVirtualMachineScaleSetID(input string) (*MaintenanceAssignmentVirtualMachineScaleSetId, error) {
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.Maintenance/configurationAssignments/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine Scale Set ID (%q)", input)
	}

	targetResourceId, name := groups[1], groups[2]
	virtualMachineScaleSetId, err := parseCompute.VirtualMachineScaleSetID(targetResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing Maintenance Assignment Virtual Machine Scale Set ID: %q: Expected valid virtual machine scale set ID", input)
	}

	return &MaintenanceAssignmentVirtualMachineScaleSetId{
		VirtualMachineScaleSetId:    virtualMachineScaleSetId,
		VirtualMachineScaleSetIdRaw: targetResourceId,
		Name:                        name,
	}, nil
}
