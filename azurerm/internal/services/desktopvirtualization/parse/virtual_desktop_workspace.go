package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualDesktopWorkspaceId struct {
	ResourceGroup string
	Name          string
}

func (id VirtualDesktopWorkspaceId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/workspaces/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func NewVirtualDesktopWorkspaceId(resourceGroup, name string) VirtualDesktopWorkspaceId {
	return VirtualDesktopWorkspaceId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func VirtualDesktopWorkspaceID(input string) (*VirtualDesktopWorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Workspace ID %q: %+v", input, err)
	}

	workspace := VirtualDesktopWorkspaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &workspace, nil
}
