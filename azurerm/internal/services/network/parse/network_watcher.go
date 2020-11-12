package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetworkWatcherId struct {
	ResourceGroup string
	Name          string
}

func NetworkWatcherID(input string) (*NetworkWatcherId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	watcher := NetworkWatcherId{
		ResourceGroup: id.ResourceGroup,
	}

	if watcher.Name, err = id.PopSegment("networkWatchers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &watcher, nil
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
