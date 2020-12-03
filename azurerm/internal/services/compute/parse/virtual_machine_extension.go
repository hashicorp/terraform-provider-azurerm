package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineExtensionId struct {
	SubscriptionId     string
	ResourceGroup      string
	VirtualMachineName string
	ExtensionName      string
}

func NewVirtualMachineExtensionID(subscriptionId, resourceGroup, virtualMachineName, extensionName string) VirtualMachineExtensionId {
	return VirtualMachineExtensionId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		VirtualMachineName: virtualMachineName,
		ExtensionName:      extensionName,
	}
}

func (id VirtualMachineExtensionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/extensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName, id.ExtensionName)
}

// VirtualMachineExtensionID parses a VirtualMachineExtension ID into an VirtualMachineExtensionId struct
func VirtualMachineExtensionID(input string) (*VirtualMachineExtensionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineExtensionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualMachineName, err = id.PopSegment("virtualMachines"); err != nil {
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
