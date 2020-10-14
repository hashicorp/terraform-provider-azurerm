package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubBgpConnectionId struct {
	ResourceGroup  string
	VirtualHubName string
	Name           string
}

func VirtualHubBgpConnectionID(input string) (*VirtualHubBgpConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing virtualHubBgpConnection ID %q: %+v", input, err)
	}

	virtualHubBgpConnection := VirtualHubBgpConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if virtualHubBgpConnection.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
		return nil, err
	}

	if virtualHubBgpConnection.Name, err = id.PopSegment("bgpConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &virtualHubBgpConnection, nil
}
