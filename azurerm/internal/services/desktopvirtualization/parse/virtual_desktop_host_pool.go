package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualDesktopHostPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func (id VirtualDesktopHostPoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostpools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func NewVirtualDesktopHostPoolId(subscriptionId, resourceGroup, name string) VirtualDesktopHostPoolId {
	return VirtualDesktopHostPoolId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

// VirtualDesktopHostPoolID - Parses and validates the virtual desktop host pool
func VirtualDesktopHostPoolID(input string) (*VirtualDesktopHostPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Host Pool ID %q: %+v", input, err)
	}

	hostPool := VirtualDesktopHostPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if hostPool.Name, err = id.PopSegment("hostpools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &hostPool, nil
}
