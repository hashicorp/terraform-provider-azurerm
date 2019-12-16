package compute

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetExtensionResourceID struct {
	ResourceGroup      string
	VirtualMachineName string
	Name               string
}

func ParseVirtualMachineScaleSetExtensionID(input string) (*VirtualMachineScaleSetExtensionResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Machine Scale Set Extension ID %q: %+v", input, err)
	}

	extension := VirtualMachineScaleSetExtensionResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if extension.VirtualMachineName, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}

	if extension.Name, err = id.PopSegment("extensions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &extension, nil
}
