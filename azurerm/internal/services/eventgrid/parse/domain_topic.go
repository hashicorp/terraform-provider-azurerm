package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DomainTopicId struct {
	ResourceGroup string
	DomainName    string
	TopicName     string
}

func DomainTopicID(input string) (*DomainTopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid Domain Topic ID %q: %+v", input, err)
	}

	domainTopic := DomainTopicId{
		ResourceGroup: id.ResourceGroup,
	}

	if domainTopic.TopicName, err = id.PopSegment("topics"); err != nil {
		return nil, err
	}

	if domainTopic.DomainName, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domainTopic, nil
}
