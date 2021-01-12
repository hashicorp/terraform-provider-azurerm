package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SystemTopicEventSubscriptionId struct {
	ResourceGroup string
	SystemTopic   string
	Name          string
}

func SystemTopicEventSubscriptionID(input string) (*SystemTopicEventSubscriptionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse EventGrid System Topic Event Subscription ID %q: %+v", input, err)
	}

	systemTopicEventSubscriptionID := SystemTopicEventSubscriptionId{
		ResourceGroup: id.ResourceGroup,
	}

	if systemTopicEventSubscriptionID.SystemTopic, err = id.PopSegment("systemTopics"); err != nil {
		return nil, err
	}

	if systemTopicEventSubscriptionID.Name, err = id.PopSegment("eventSubscriptions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &systemTopicEventSubscriptionID, nil
}
