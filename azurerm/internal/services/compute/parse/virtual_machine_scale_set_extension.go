package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineScaleSetExtensionId struct {
	SubscriptionId             string
	ResourceGroup              string
	VirtualMachineScaleSetName string
	ExtensionName              string
}

func NewVirtualMachineScaleSetExtensionID(subscriptionId, resourceGroup, virtualMachineScaleSetName, extensionName string) VirtualMachineScaleSetExtensionId {
	return VirtualMachineScaleSetExtensionId{
		SubscriptionId:             subscriptionId,
		ResourceGroup:              resourceGroup,
		VirtualMachineScaleSetName: virtualMachineScaleSetName,
		ExtensionName:              extensionName,
	}
}

func (id VirtualMachineScaleSetExtensionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s/extensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineScaleSetName, id.ExtensionName)
}

func VirtualMachineScaleSetExtensionID(input string) (*VirtualMachineScaleSetExtensionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineScaleSetExtensionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.VirtualMachineScaleSetName, err = id.PopSegment("virtualMachineScaleSets"); err != nil {
		return nil, err
	}
	if resourceId.ExtensionName, err = id.PopSegment("extensions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
