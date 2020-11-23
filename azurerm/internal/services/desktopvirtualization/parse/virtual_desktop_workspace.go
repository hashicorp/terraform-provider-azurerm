package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualDesktopWorkspaceId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func (id VirtualDesktopWorkspaceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/workspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func NewVirtualDesktopWorkspaceId(subscriptionId, resourceGroup, name string) VirtualDesktopWorkspaceId {
	return VirtualDesktopWorkspaceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func VirtualDesktopWorkspaceID(input string) (*VirtualDesktopWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Workspace ID %q: %+v", input, err)
	}

	workspace := VirtualDesktopWorkspaceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &workspace, nil
}
