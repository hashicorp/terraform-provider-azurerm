package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventGridTopicId struct {
	ResourceGroup string
	Name          string
}

func EventGridTopicID(input string) (*EventGridTopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid Topic ID %q: %+v", input, err)
	}

	topic := EventGridTopicId{
		ResourceGroup: id.ResourceGroup,
	}

	if topic.Name, err = id.PopSegment("topics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &topic, nil
}
