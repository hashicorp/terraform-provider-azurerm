package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventGridDomainTopicId struct {
	ResourceGroup string
	Name          string
	Domain        string
}

func EventGridDomainTopicID(input string) (*EventGridDomainTopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid Domain Topic ID %q: %+v", input, err)
	}

	domainTopic := EventGridDomainTopicId{
		ResourceGroup: id.ResourceGroup,
	}

	if domainTopic.Name, err = id.PopSegment("topics"); err != nil {
		return nil, err
	}

	if domainTopic.Domain, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domainTopic, nil
}
