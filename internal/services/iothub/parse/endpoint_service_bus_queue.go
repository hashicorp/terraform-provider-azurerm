package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type EndpointServiceBusQueueId struct {
	SubscriptionId string
	ResourceGroup  string
	IotHubName     string
	EndpointName   string
}

func NewEndpointServiceBusQueueID(subscriptionId, resourceGroup, iotHubName, endpointName string) EndpointServiceBusQueueId {
	return EndpointServiceBusQueueId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		IotHubName:     iotHubName,
		EndpointName:   endpointName,
	}
}

func (id EndpointServiceBusQueueId) String() string {
	segments := []string{
		fmt.Sprintf("Endpoint Name %q", id.EndpointName),
		fmt.Sprintf("Iot Hub Name %q", id.IotHubName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Endpoint Service Bus Queue", segmentsStr)
}

func (id EndpointServiceBusQueueId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/IotHubs/%s/Endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IotHubName, id.EndpointName)
}

// EndpointServiceBusQueueID parses a EndpointServiceBusQueue ID into an EndpointServiceBusQueueId struct
func EndpointServiceBusQueueID(input string) (*EndpointServiceBusQueueId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EndpointServiceBusQueueId{
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
	if resourceId.EndpointName, err = id.PopSegment("Endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
