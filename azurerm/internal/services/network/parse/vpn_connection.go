package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VpnConnectionId struct {
	ResourceGroup  string
	VpnGatewayName string
	Name           string
}

func (id VpnConnectionId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/vpnConnections/%s",
		subscriptionId, id.ResourceGroup, id.VpnGatewayName, id.Name)
}

func NewVpnConnectionID(resourceGroup, gateway, name string) VpnConnectionId {
	return VpnConnectionId{
		ResourceGroup:  resourceGroup,
		VpnGatewayName: gateway,
		Name:           name,
	}
}

func VpnConnectionID(input string) (*VpnConnectionId, error) {
	rawId, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing VPNGateway Connection ID %q: %+v", input, err)
	}

	id := VpnConnectionId{
		ResourceGroup: rawId.ResourceGroup,
	}

	if id.VpnGatewayName, err = rawId.PopSegment("vpnGateways"); err != nil {
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
