package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetExtensionId struct {
	ResourceGroup              string
	VirtualMachineScaleSetName string
	Name                       string
}

func NewVirtualMachineScaleSetExtensionId(id VirtualMachineScaleSetId, name string) VirtualMachineScaleSetExtensionId {
	return VirtualMachineScaleSetExtensionId{
		ResourceGroup:              id.ResourceGroup,
		VirtualMachineScaleSetName: id.Name,
		Name:                       name,
	}
}

func (id VirtualMachineScaleSetExtensionId) ID(subscriptionId string) string {
	base := NewVirtualMachineScaleSetId(id.ResourceGroup, id.VirtualMachineScaleSetName).ID(subscriptionId)
	return fmt.Sprintf("%s/extensions/%s", base, id.Name)
}

func VirtualMachineScaleSetExtensionID(input string) (*VirtualMachineScaleSetExtensionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Virtual Machine Scale Set Extension ID %q: %+v", input, err)
	}

	extension := VirtualMachineScaleSetExtensionId{
		ResourceGroup: id.ResourceGroup,
	}

	if extension.VirtualMachineScaleSetName, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
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
