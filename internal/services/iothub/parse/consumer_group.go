package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ConsumerGroupId struct {
	SubscriptionId       string
	ResourceGroup        string
	IotHubName           string
	EventHubEndpointName string
	Name                 string
}

func NewConsumerGroupID(subscriptionId, resourceGroup, iotHubName, eventHubEndpointName, name string) ConsumerGroupId {
	return ConsumerGroupId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		IotHubName:           iotHubName,
		EventHubEndpointName: eventHubEndpointName,
		Name:                 name,
	}
}

func (id ConsumerGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Event Hub Endpoint Name %q", id.EventHubEndpointName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Consumer Group", segmentsStr)
}

func (id ConsumerGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/IotHubs/%s/eventHubEndpoints/%s/ConsumerGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.EventHubEndpointName, id.Name)
}

// ConsumerGroupID parses a ConsumerGroup ID into an ConsumerGroupId struct
func ConsumerGroupID(input string) (*ConsumerGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConsumerGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IotHubName, err = id.PopSegment("IotHubs"); err != nil {
		return nil, err
	}
	if resourceId.EventHubEndpointName, err = id.PopSegment("eventHubEndpoints"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("ConsumerGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
