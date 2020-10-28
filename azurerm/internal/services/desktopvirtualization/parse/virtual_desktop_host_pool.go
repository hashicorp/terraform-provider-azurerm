package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualDesktopHostPoolId struct {
	ResourceGroup string
	Name          string
}

func (id VirtualDesktopHostPoolId) ID(subscriptionId string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostpools/%s"
	return fmt.Sprintf(fmtString, subscriptionId, id.ResourceGroup, id.Name)
}

func NewVirtualDesktopHostPoolId(resourceGroup, name string) VirtualDesktopHostPoolId {
	return VirtualDesktopHostPoolId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

// VirtualDesktopHostPoolID - Parses and validates the virtual desktop host pool
func VirtualDesktopHostPoolID(input string) (*VirtualDesktopHostPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Host Pool ID %q: %+v", input, err)
	}

	hostPool := VirtualDesktopHostPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if hostPool.Name, err = id.PopSegment("hostpools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &hostPool, nil
}
