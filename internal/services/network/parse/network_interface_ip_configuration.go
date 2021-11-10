package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type NetworkInterfaceIpConfigurationId struct {
	SubscriptionId       string
	ResourceGroup        string
	NetworkInterfaceName string
	IpConfigurationName  string
}

func NewNetworkInterfaceIpConfigurationID(subscriptionId, resourceGroup, networkInterfaceName, ipConfigurationName string) NetworkInterfaceIpConfigurationId {
	return NetworkInterfaceIpConfigurationId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		NetworkInterfaceName: networkInterfaceName,
		IpConfigurationName:  ipConfigurationName,
	}
}

func (id NetworkInterfaceIpConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Ip Configuration Name %q", id.IpConfigurationName),
		fmt.Sprintf("Network Interface Name %q", id.NetworkInterfaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Interface Ip Configuration", segmentsStr)
}

func (id NetworkInterfaceIpConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkInterfaceName, id.IpConfigurationName)
}

// NetworkInterfaceIpConfigurationID parses a NetworkInterfaceIpConfiguration ID into an NetworkInterfaceIpConfigurationId struct
func NetworkInterfaceIpConfigurationID(input string) (*NetworkInterfaceIpConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkInterfaceIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkInterfaceName, err = id.PopSegment("networkInterfaces"); err != nil {
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
