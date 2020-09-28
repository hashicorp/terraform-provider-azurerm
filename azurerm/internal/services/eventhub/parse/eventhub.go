package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type EventHubId struct {
	ResourceGroup string
	Namespace     string
	Name          string
}

func EventHubID(input string) (*EventHubId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	eventHubId := EventHubId{
		ResourceGroup: id.ResourceGroup,
	}

	if eventHubId.Namespace, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}

	if eventHubId.Name, err = id.PopSegment("eventhubs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &eventHubId, nil
}
