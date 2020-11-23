package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualDesktopApplicationGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewVirtualDesktopApplicationGroupId(subscriptionId, resourceGroup, name string) VirtualDesktopApplicationGroupId {
	return VirtualDesktopApplicationGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id VirtualDesktopApplicationGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/applicationgroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func VirtualDesktopApplicationGroupID(input string) (*VirtualDesktopApplicationGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Application Group ID %q: %+v", input, err)
	}

	applicationGroup := VirtualDesktopApplicationGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if applicationGroup.Name, err = id.PopSegment("applicationgroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &applicationGroup, nil
}
