package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PointToSiteVPNGatewayResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParsePointToSiteVPNGatewayID(input string) (*PointToSiteVPNGatewayResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Point-to-Site VPN Gateway ID %q: %+v", input, err)
	}

	routeTable := PointToSiteVPNGatewayResourceID{
		Base: *id,
		Name: id.Path["p2svpnGateways"],
	}

	if routeTable.Name == "" {
		return nil, fmt.Errorf("ID was missing the `p2svpnGateways` element")
	}

	return &routeTable, nil
}
