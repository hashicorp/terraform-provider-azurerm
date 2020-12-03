package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DigitalTwinsEndpointId struct {
	SubscriptionId           string
	ResourceGroup            string
	DigitalTwinsInstanceName string
	EndpointName             string
}

func NewDigitalTwinsEndpointID(subscriptionId, resourceGroup, digitalTwinsInstanceName, endpointName string) DigitalTwinsEndpointId {
	return DigitalTwinsEndpointId{
		SubscriptionId:           subscriptionId,
		ResourceGroup:            resourceGroup,
		DigitalTwinsInstanceName: digitalTwinsInstanceName,
		EndpointName:             endpointName,
	}
}

func (id DigitalTwinsEndpointId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s/endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DigitalTwinsInstanceName, id.EndpointName)
}

func DigitalTwinsEndpointID(input string) (*DigitalTwinsEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DigitalTwinsEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DigitalTwinsInstanceName, err = id.PopSegment("digitalTwinsInstances"); err != nil {
		return nil, err
	}
	if resourceId.EndpointName, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
