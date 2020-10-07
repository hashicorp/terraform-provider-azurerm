package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetId struct {
	ResourceGroup string
	Name          string
}

func NewVirtualMachineScaleSetId(resourceGroup, name string) VirtualMachineScaleSetId {
	return VirtualMachineScaleSetId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id VirtualMachineScaleSetId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func VirtualMachineScaleSetID(input string) (*VirtualMachineScaleSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Virtual Machine Scale Set ID %q: %+v", input, err)
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
