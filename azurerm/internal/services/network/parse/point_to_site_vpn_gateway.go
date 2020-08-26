package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PointToSiteVPNGatewayId struct {
	ResourceGroup string
	Name          string
}

func PointToSiteVPNGatewayID(input string) (*PointToSiteVPNGatewayId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Point-to-Site VPN Gateway ID %q: %+v", input, err)
	}

	routeTable := PointToSiteVPNGatewayId{
		ResourceGroup: id.ResourceGroup,
	}

	if routeTable.Name, err = id.PopSegment("p2sVpnGateways"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &routeTable, nil
}
