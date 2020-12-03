package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VpnConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	VpnGatewayName string
	Name           string
}

func NewVpnConnectionID(subscriptionId, resourceGroup, vpnGatewayName, name string) VpnConnectionId {
	return VpnConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VpnGatewayName: vpnGatewayName,
		Name:           name,
	}
}

func (id VpnConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/vpnConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VpnGatewayName, id.Name)
}

// VpnConnectionID parses a VpnConnection ID into an VpnConnectionId struct
func VpnConnectionID(input string) (*VpnConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VpnConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VpnGatewayName, err = id.PopSegment("vpnGateways"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("vpnConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
