package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkManagerSubscriptionConnectionId struct {
	SubscriptionId               string
	NetworkManagerConnectionName string
}

func NewNetworkManagerSubscriptionConnectionID(subscriptionId, networkManagerConnectionName string) NetworkManagerSubscriptionConnectionId {
	return NetworkManagerSubscriptionConnectionId{
		SubscriptionId:               subscriptionId,
		NetworkManagerConnectionName: networkManagerConnectionName,
	}
}

func (id NetworkManagerSubscriptionConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Network Manager Connection Name %q", id.NetworkManagerConnectionName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Subscription Connection", segmentsStr)
}

func (id NetworkManagerSubscriptionConnectionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/networkManagerConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.NetworkManagerConnectionName)
}

// NetworkManagerSubscriptionConnectionID parses a NetworkManagerSubscriptionConnection ID into an NetworkManagerSubscriptionConnectionId struct
func NetworkManagerSubscriptionConnectionID(input string) (*NetworkManagerSubscriptionConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerSubscriptionConnectionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.NetworkManagerConnectionName, err = id.PopSegment("networkManagerConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
