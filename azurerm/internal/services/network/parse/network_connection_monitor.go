package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkConnectionMonitorId struct {
	ResourceGroup string
	WatcherName   string
	Name          string
}

func NetworkConnectionMonitorID(input string) (*NetworkConnectionMonitorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	connectionMonitor := NetworkConnectionMonitorId{
		ResourceGroup: id.ResourceGroup,
	}

	if connectionMonitor.WatcherName, err = id.PopSegment("networkWatchers"); err != nil {
		return nil, err
	}

	if connectionMonitor.Name, err = id.PopSegment("connectionMonitors"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &connectionMonitor, nil
}

func (id NetworkWatcherId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s",
		subscriptionId, id.ResourceGroup, id.Name)
}

func NewNetworkWatcherID(resourceGroup, name string) NetworkWatcherId {
	return NetworkWatcherId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}
