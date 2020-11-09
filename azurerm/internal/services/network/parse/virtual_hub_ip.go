package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubIPId struct {
	ResourceGroup  string
	VirtualHubName string
	Name           string
}

func VirtualHubIPID(input string) (*VirtualHubIPId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing virtualHubIP ID %q: %+v", input, err)
	}

	virtualHubIP := VirtualHubIPId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHubIP.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if virtualHubIP.Name, err = id.PopSegment("ipConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHubIP, nil
}
