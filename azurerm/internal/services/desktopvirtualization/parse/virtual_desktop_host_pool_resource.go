package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// VirtualDesktopHostPoolid - The id for the virtual desktop host pool
type VirtualDesktopHostPoolid struct {
	ResourceGroup string
	Name          string
}

// VirtualDesktopHostPoolID - Parses and validates the virtual desktop host pool
func VirtualDesktopHostPoolID(input string) (*VirtualDesktopHostPoolid, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Host Pool ID %q: %+v", input, err)
	}

	hostPool := VirtualDesktopHostPoolid{
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
