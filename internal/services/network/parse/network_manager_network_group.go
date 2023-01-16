package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkManagerNetworkGroupId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkManagerName string
	NetworkGroupName   string
}

func NewNetworkManagerNetworkGroupID(subscriptionId, resourceGroup, networkManagerName, networkGroupName string) NetworkManagerNetworkGroupId {
	return NetworkManagerNetworkGroupId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkManagerName: networkManagerName,
		NetworkGroupName:   networkGroupName,
	}
}

func (id NetworkManagerNetworkGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Network Group Name %q", id.NetworkGroupName),
		fmt.Sprintf("Network Manager Name %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Network Group", segmentsStr)
}

func (id NetworkManagerNetworkGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/networkGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName)
}

// NetworkManagerNetworkGroupID parses a NetworkManagerNetworkGroup ID into an NetworkManagerNetworkGroupId struct
func NetworkManagerNetworkGroupID(input string) (*NetworkManagerNetworkGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerNetworkGroupId{
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
	if resourceId.NetworkGroupName, err = id.PopSegment("networkGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
