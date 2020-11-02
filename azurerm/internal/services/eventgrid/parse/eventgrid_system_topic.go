package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventGridSystemTopicId struct {
	ResourceGroup string
	Name          string
}

func EventGridSystemTopicID(input string) (*EventGridSystemTopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid System Topic ID %q: %+v", input, err)
	}

	topic := EventGridSystemTopicId{
		ResourceGroup: id.ResourceGroup,
	}

	if topic.Name, err = id.PopSegment("systemTopics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &topic, nil
}
