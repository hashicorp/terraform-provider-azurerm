package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkManagerStaticMemberId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkManagerName string
	NetworkGroupName   string
	StaticMemberName   string
}

func NewNetworkManagerStaticMemberID(subscriptionId, resourceGroup, networkManagerName, networkGroupName, staticMemberName string) NetworkManagerStaticMemberId {
	return NetworkManagerStaticMemberId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkManagerName: networkManagerName,
		NetworkGroupName:   networkGroupName,
		StaticMemberName:   staticMemberName,
	}
}

func (id NetworkManagerStaticMemberId) String() string {
	segments := []string{
		fmt.Sprintf("Static Member Name %q", id.StaticMemberName),
		fmt.Sprintf("Network Group Name %q", id.NetworkGroupName),
		fmt.Sprintf("Network Manager Name %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Static Member", segmentsStr)
}

func (id NetworkManagerStaticMemberId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/networkGroups/%s/staticMembers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.NetworkGroupName, id.StaticMemberName)
}

// NetworkManagerStaticMemberID parses a NetworkManagerStaticMember ID into an NetworkManagerStaticMemberId struct
func NetworkManagerStaticMemberID(input string) (*NetworkManagerStaticMemberId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerStaticMemberId{
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
	if resourceId.StaticMemberName, err = id.PopSegment("staticMembers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
