package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewVirtualMachineId(subscriptionId, resourceGroup, name string) VirtualMachineId {
	return VirtualMachineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id VirtualMachineId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func VirtualMachineID(input string) (*VirtualMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Virtual Machine ID %q: %+v", input, err)
	}

	vm := VirtualMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if vm.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vm, nil
}
