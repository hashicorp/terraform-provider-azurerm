package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SystemTopicEventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroup         string
	SystemTopicName       string
	EventSubscriptionName string
}

func NewSystemTopicEventSubscriptionID(subscriptionId, resourceGroup, systemTopicName, eventSubscriptionName string) SystemTopicEventSubscriptionId {
	return SystemTopicEventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		SystemTopicName:       systemTopicName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

func (id SystemTopicEventSubscriptionId) String() string {
	segments := []string{
		fmt.Sprintf("Event Subscription Name %q", id.EventSubscriptionName),
		fmt.Sprintf("System Topic Name %q", id.SystemTopicName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "System Topic Event Subscription", segmentsStr)
}

func (id SystemTopicEventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/systemTopics/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SystemTopicName, id.EventSubscriptionName)
}

// SystemTopicEventSubscriptionID parses a SystemTopicEventSubscription ID into an SystemTopicEventSubscriptionId struct
func SystemTopicEventSubscriptionID(input string) (*SystemTopicEventSubscriptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SystemTopicEventSubscriptionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SystemTopicName, err = id.PopSegment("systemTopics"); err != nil {
		return nil, err
	}
	if resourceId.EventSubscriptionName, err = id.PopSegment("eventSubscriptions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
