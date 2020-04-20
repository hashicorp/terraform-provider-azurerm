package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetId struct {
	ResourceGroup string
	Name          string
}

func VirtualMachineScaleSetID(input string) (*VirtualMachineScaleSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Machine Scale Set ID %q: %+v", input, err)
	}

	vmScaleSet := VirtualMachineScaleSetId{
		ResourceGroup: id.ResourceGroup,
	}

	if vmScaleSet.Name, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vmScaleSet, nil
}
