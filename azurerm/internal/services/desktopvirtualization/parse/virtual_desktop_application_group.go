package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

// VirtualDesktopApplicationGroupid - The id for the virtual desktop host pool
type VirtualDesktopApplicationGroupid struct {
	ResourceGroup string
	Name          string
}

// VirtualDesktopApplicationGroupID - Parses and validates the virtual desktop host pool
func VirtualDesktopApplicationGroupID(input string) (*VirtualDesktopApplicationGroupid, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Virtual Desktop Application Group ID %q: %+v", input, err)
	}

	ApplicationGroup := VirtualDesktopApplicationGroupid{
		ResourceGroup: id.ResourceGroup,
	}

	if ApplicationGroup.Name, err = id.PopSegment("applicationgroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &ApplicationGroup, nil
}
