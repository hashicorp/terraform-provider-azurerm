package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VPNGatewayConnectionId struct {
	ResourceGroup string
	Gateway       string
	Name          string
}

func (id VPNGatewayConnectionId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/vpnConnections/%s",
		subscriptionId, id.ResourceGroup, id.Gateway, id.Name)
}

func NewVPNGatewayConnectionID(resourceGroup, gateway, name string) VPNGatewayConnectionId {
	return VPNGatewayConnectionId{
		ResourceGroup: resourceGroup,
		Gateway:       gateway,
		Name:          name,
	}
}

func VPNGatewayConnectionID(input string) (*VPNGatewayConnectionId, error) {
	rawId, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing VPNGateway Connection ID %q: %+v", input, err)
	}

	id := VPNGatewayConnectionId{
		ResourceGroup: rawId.ResourceGroup,
	}

	if id.Gateway, err = rawId.PopSegment("vpnGateways"); err != nil {
		return nil, err
	}

	if id.Name, err = rawId.PopSegment("vpnConnections"); err != nil {
		return nil, err
	}

	if err := rawId.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &id, nil
}

type VPNGatewayId struct {
	ResourceGroup string
	Name          string
}

func (id VPNGatewayId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewVPNGatewayID(resourceGroup, name string) VPNGatewayId {
	return VPNGatewayId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func VPNGatewayID(input string) (*VPNGatewayId, error) {
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
