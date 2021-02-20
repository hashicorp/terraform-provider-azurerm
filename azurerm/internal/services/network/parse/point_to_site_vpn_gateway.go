package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PointToSiteVpnGatewayId struct {
	SubscriptionId    string
	ResourceGroup     string
	P2sVpnGatewayName string
}

func NewPointToSiteVpnGatewayID(subscriptionId, resourceGroup, p2sVpnGatewayName string) PointToSiteVpnGatewayId {
	return PointToSiteVpnGatewayId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		P2sVpnGatewayName: p2sVpnGatewayName,
	}
}

func (id PointToSiteVpnGatewayId) String() string {
	segments := []string{
		fmt.Sprintf("P2s Vpn Gateway Name %q", id.P2sVpnGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Point To Site Vpn Gateway", segmentsStr)
}

func (id PointToSiteVpnGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/p2sVpnGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.P2sVpnGatewayName)
}

// PointToSiteVpnGatewayID parses a PointToSiteVpnGateway ID into an PointToSiteVpnGatewayId struct
func PointToSiteVpnGatewayID(input string) (*PointToSiteVpnGatewayId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PointToSiteVpnGatewayId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.P2sVpnGatewayName, err = id.PopSegment("p2sVpnGateways"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
