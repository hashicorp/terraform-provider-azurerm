package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ConnectionMonitorId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkWatcherName string
	Name               string
}

func NewConnectionMonitorID(subscriptionId, resourceGroup, networkWatcherName, name string) ConnectionMonitorId {
	return ConnectionMonitorId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkWatcherName: networkWatcherName,
		Name:               name,
	}
}

func (id ConnectionMonitorId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Network Watcher Name %q", id.NetworkWatcherName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Connection Monitor", segmentsStr)
}

func (id ConnectionMonitorId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/connectionMonitors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkWatcherName, id.Name)
}

// ConnectionMonitorID parses a ConnectionMonitor ID into an ConnectionMonitorId struct
func ConnectionMonitorID(input string) (*ConnectionMonitorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConnectionMonitorId{
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
	if resourceId.Name, err = id.PopSegment("connectionMonitors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
