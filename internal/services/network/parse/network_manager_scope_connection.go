package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkManagerScopeConnectionId struct {
	SubscriptionId      string
	ResourceGroup       string
	NetworkManagerName  string
	ScopeConnectionName string
}

func NewNetworkManagerScopeConnectionID(subscriptionId, resourceGroup, networkManagerName, scopeConnectionName string) NetworkManagerScopeConnectionId {
	return NetworkManagerScopeConnectionId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		NetworkManagerName:  networkManagerName,
		ScopeConnectionName: scopeConnectionName,
	}
}

func (id NetworkManagerScopeConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Scope Connection Name %q", id.ScopeConnectionName),
		fmt.Sprintf("Network Manager Name %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Scope Connection", segmentsStr)
}

func (id NetworkManagerScopeConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/scopeConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.ScopeConnectionName)
}

// NetworkManagerScopeConnectionID parses a NetworkManagerScopeConnection ID into an NetworkManagerScopeConnectionId struct
func NetworkManagerScopeConnectionID(input string) (*NetworkManagerScopeConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerScopeConnectionId{
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
	if resourceId.ScopeConnectionName, err = id.PopSegment("scopeConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
