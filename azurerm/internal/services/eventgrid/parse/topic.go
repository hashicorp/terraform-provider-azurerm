package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TopicId struct {
	ResourceGroup string
	Name          string
}

func TopicID(input string) (*TopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid Topic ID %q: %+v", input, err)
	}

	topic := TopicId{
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
