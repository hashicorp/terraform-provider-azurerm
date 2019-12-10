package compute

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetExtensionResourceID struct {
	Base azure.ResourceID

	VirtualMachineName string
	Name               string
}

func ParseVirtualMachineScaleSetExtensionResourceID(input string) (*VirtualMachineScaleSetExtensionResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Machine Scale Set Extension ID %q: %+v", input, err)
	}

	extension := VirtualMachineScaleSetExtensionResourceID{
		Base:               *id,
		VirtualMachineName: id.Path["virtualMachineScaleSets"],
		Name:               id.Path["extensions"],
	}

	if extension.VirtualMachineName == "" {
		return nil, fmt.Errorf("ID was missing the `virtualMachineScaleSets` element")
	}

	if extension.Name == "" {
		return nil, fmt.Errorf("ID was missing the `extensions` element")
	}

	return &extension, nil
}
