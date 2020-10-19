package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineId struct {
	ResourceGroup string
	Name          string
}

func NewVirtualMachineId(resourceGroup, name string) VirtualMachineId {
	return VirtualMachineId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id VirtualMachineId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func VirtualMachineID(input string) (*VirtualMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Virtual Machine ID %q: %+v", input, err)
	}

	vm := VirtualMachineId{
		ResourceGroup: id.ResourceGroup,
	}

	if vm.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vm, nil
}
