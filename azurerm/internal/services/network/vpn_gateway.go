package network

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VPNGatewayResourceID struct {
	Base azure.ResourceID

	Name string
}

func ParseVPNGatewayID(input string) (*VPNGatewayResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse VPN Gateway ID %q: %+v", input, err)
	}

	gateway := VPNGatewayResourceID{
		Base: *id,
		Name: id.Path["vpnGateways"],
	}

	if gateway.Name == "" {
		return nil, fmt.Errorf("ID was missing the `vpnGateways` element")
	}

	return &gateway, nil
}
