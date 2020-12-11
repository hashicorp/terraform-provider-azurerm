package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id DigitalTwinsEndpointId) String() string {
	segments := []string{
		fmt.Sprintf("Endpoint Name %q", id.EndpointName),
		fmt.Sprintf("Digital Twins Instance Name %q", id.DigitalTwinsInstanceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Digital Twins Endpoint", segmentsStr)
}

func (id DigitalTwinsEndpointId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s/endpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DigitalTwinsInstanceName, id.EndpointName)
}

// DigitalTwinsEndpointID parses a DigitalTwinsEndpoint ID into an DigitalTwinsEndpointId struct
func DigitalTwinsEndpointID(input string) (*DigitalTwinsEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DigitalTwinsEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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
