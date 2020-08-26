package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VPNGatewayId struct {
	ResourceGroup string
	Name          string
}

func ParseVPNGatewayID(input string) (*VPNGatewayId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse VPN Gateway ID %q: %+v", input, err)
	}

	gateway := VPNGatewayId{
		ResourceGroup: id.ResourceGroup,
	}

	if gateway.Name, err = id.PopSegment("vpnGateways"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &gateway, nil
}
