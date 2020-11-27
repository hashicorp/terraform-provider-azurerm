package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ConnectionMonitorId struct {
	ResourceGroup      string
	NetworkWatcherName string
	Name               string
}

func ConnectionMonitorID(input string) (*ConnectionMonitorId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	connectionMonitor := ConnectionMonitorId{
		ResourceGroup: id.ResourceGroup,
	}

	if connectionMonitor.NetworkWatcherName, err = id.PopSegment("networkWatchers"); err != nil {
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
