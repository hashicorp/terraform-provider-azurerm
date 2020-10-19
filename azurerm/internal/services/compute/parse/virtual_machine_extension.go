package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineExtensionId struct {
	ResourceGroup  string
	Name           string
	VirtualMachine string
}

func NewVirtualMachineExtensionId(id VirtualMachineId, name string) VirtualMachineExtensionId {
	return VirtualMachineExtensionId{
		ResourceGroup:  id.ResourceGroup,
		VirtualMachine: id.Name,
		Name:           name,
	}
}

func (id VirtualMachineExtensionId) ID(subscriptionId string) string {
	base := NewVirtualMachineId(id.ResourceGroup, id.VirtualMachine).ID(subscriptionId)
	return fmt.Sprintf("%s/extensions/%s", base, id.Name)
}

func VirtualMachineExtensionID(input string) (*VirtualMachineExtensionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Virtual Machine Extension ID %q: %+v", input, err)
	}

	virtualMachineExtension := VirtualMachineExtensionId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualMachineExtension.VirtualMachine, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if virtualMachineExtension.Name, err = id.PopSegment("extensions"); err != nil {
		return nil, err
	}

	return &virtualMachineExtension, nil
}
