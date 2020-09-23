package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

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
