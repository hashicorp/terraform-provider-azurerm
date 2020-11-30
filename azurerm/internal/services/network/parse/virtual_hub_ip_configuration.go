package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualHubIpConfigurationId struct {
	SubscriptionId      string
	ResourceGroup       string
	VirtualHubName      string
	IpConfigurationName string
}

func NewVirtualHubIpConfigurationID(subscriptionId, resourceGroup, virtualHubName, ipConfigurationName string) VirtualHubIpConfigurationId {
	return VirtualHubIpConfigurationId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		VirtualHubName:      virtualHubName,
		IpConfigurationName: ipConfigurationName,
	}
}

func (id VirtualHubIpConfigurationId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.IpConfigurationName)
}

// VirtualHubIpConfigurationID parses a VirtualHubIpConfiguration ID into an VirtualHubIpConfigurationId struct
func VirtualHubIpConfigurationID(input string) (*VirtualHubIpConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualHubIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualHubName, err = id.PopSegment("virtualHubs"); err != nil {
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
