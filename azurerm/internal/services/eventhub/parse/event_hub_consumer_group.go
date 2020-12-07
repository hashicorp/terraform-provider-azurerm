package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EventHubConsumerGroupId struct {
	SubscriptionId    string
	ResourceGroup     string
	NamespaceName     string
	EventhubName      string
	ConsumergroupName string
}

func NewEventHubConsumerGroupID(subscriptionId, resourceGroup, namespaceName, eventhubName, consumergroupName string) EventHubConsumerGroupId {
	return EventHubConsumerGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		NamespaceName:     namespaceName,
		EventhubName:      eventhubName,
		ConsumergroupName: consumergroupName,
	}
}

func (id EventHubConsumerGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s/consumergroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.EventhubName, id.ConsumergroupName)
}

// EventHubConsumerGroupID parses a EventHubConsumerGroup ID into an EventHubConsumerGroupId struct
func EventHubConsumerGroupID(input string) (*EventHubConsumerGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EventHubConsumerGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.EventhubName, err = id.PopSegment("eventhubs"); err != nil {
		return nil, err
	}
	if resourceId.ConsumergroupName, err = id.PopSegment("consumergroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
