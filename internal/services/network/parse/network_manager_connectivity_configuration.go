package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkManagerConnectivityConfigurationId struct {
	SubscriptionId                string
	ResourceGroup                 string
	NetworkManagerName            string
	ConnectivityConfigurationName string
}

func NewNetworkManagerConnectivityConfigurationID(subscriptionId, resourceGroup, networkManagerName, connectivityConfigurationName string) NetworkManagerConnectivityConfigurationId {
	return NetworkManagerConnectivityConfigurationId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		NetworkManagerName:            networkManagerName,
		ConnectivityConfigurationName: connectivityConfigurationName,
	}
}

func (id NetworkManagerConnectivityConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Connectivity Configuration Name %q", id.ConnectivityConfigurationName),
		fmt.Sprintf("Network Manager Name %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Connectivity Configuration", segmentsStr)
}

func (id NetworkManagerConnectivityConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/connectivityConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
}

// NetworkManagerConnectivityConfigurationID parses a NetworkManagerConnectivityConfiguration ID into an NetworkManagerConnectivityConfigurationId struct
func NetworkManagerConnectivityConfigurationID(input string) (*NetworkManagerConnectivityConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerConnectivityConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkManagerName, err = id.PopSegment("networkManagers"); err != nil {
		return nil, err
	}
	if resourceId.ConnectivityConfigurationName, err = id.PopSegment("connectivityConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
