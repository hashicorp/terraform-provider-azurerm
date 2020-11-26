package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DomainTopicId struct {
	SubscriptionId string
	ResourceGroup  string
	DomainName     string
	TopicName      string
}

func NewDomainTopicID(subscriptionId, resourceGroup, domainName, topicName string) DomainTopicId {
	return DomainTopicId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		DomainName:     domainName,
		TopicName:      topicName,
	}
}

func (id DomainTopicId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/domains/%s/topics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DomainName, id.TopicName)
}

// DomainTopicID parses a DomainTopic ID into an DomainTopicId struct
func DomainTopicID(input string) (*DomainTopicId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DomainTopicId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DomainName, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}
	if resourceId.TopicName, err = id.PopSegment("topics"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
