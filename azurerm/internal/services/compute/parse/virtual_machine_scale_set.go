package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewVirtualMachineScaleSetId(subscriptionId, resourceGroup, name string) VirtualMachineScaleSetId {
	return VirtualMachineScaleSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id VirtualMachineScaleSetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func VirtualMachineScaleSetID(input string) (*VirtualMachineScaleSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Virtual Machine Scale Set ID %q: %+v", input, err)
	}

	vmScaleSet := VirtualMachineScaleSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if vmScaleSet.Name, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vmScaleSet, nil
}
