package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualNetworkGatewayIpConfigurationId struct {
	SubscriptionId            string
	ResourceGroup             string
	VirtualNetworkGatewayName string
	IpConfigurationName       string
}

func NewVirtualNetworkGatewayIpConfigurationID(subscriptionId, resourceGroup, virtualNetworkGatewayName, ipConfigurationName string) VirtualNetworkGatewayIpConfigurationId {
	return VirtualNetworkGatewayIpConfigurationId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		VirtualNetworkGatewayName: virtualNetworkGatewayName,
		IpConfigurationName:       ipConfigurationName,
	}
}

func (id VirtualNetworkGatewayIpConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Ip Configuration Name %q", id.IpConfigurationName),
		fmt.Sprintf("Virtual Network Gateway Name %q", id.VirtualNetworkGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Network Gateway Ip Configuration", segmentsStr)
}

func (id VirtualNetworkGatewayIpConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworkGateways/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkGatewayName, id.IpConfigurationName)
}

// VirtualNetworkGatewayIpConfigurationID parses a VirtualNetworkGatewayIpConfiguration ID into an VirtualNetworkGatewayIpConfigurationId struct
func VirtualNetworkGatewayIpConfigurationID(input string) (*VirtualNetworkGatewayIpConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualNetworkGatewayIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualNetworkGatewayName, err = id.PopSegment("virtualNetworkGateways"); err != nil {
		return nil, err
	}
	if resourceId.IpConfigurationName, err = id.PopSegment("ipConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
