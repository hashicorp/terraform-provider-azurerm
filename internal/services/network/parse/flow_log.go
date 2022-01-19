package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FlowLogId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkWatcherName string
	Name               string
}

func NewFlowLogID(subscriptionId, resourceGroup, networkWatcherName, name string) FlowLogId {
	return FlowLogId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkWatcherName: networkWatcherName,
		Name:               name,
	}
}

func (id FlowLogId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Network Watcher Name %q", id.NetworkWatcherName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Flow Log", segmentsStr)
}

func (id FlowLogId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/flowLogs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkWatcherName, id.Name)
}

// FlowLogID parses a FlowLog ID into an FlowLogId struct
func FlowLogID(input string) (*FlowLogId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FlowLogId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkWatcherName, err = id.PopSegment("networkWatchers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("flowLogs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
