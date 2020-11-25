package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SystemTopicId struct {
	ResourceGroup string
	Name          string
}

func SystemTopicID(input string) (*SystemTopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid System Topic ID %q: %+v", input, err)
	}

	topic := SystemTopicId{
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
