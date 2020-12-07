package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AzureEndpointId struct {
	SubscriptionId            string
	ResourceGroup             string
	TrafficManagerProfileName string
	Name                      string
}

func NewAzureEndpointID(subscriptionId, resourceGroup, trafficManagerProfileName, name string) AzureEndpointId {
	return AzureEndpointId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		TrafficManagerProfileName: trafficManagerProfileName,
		Name:                      name,
	}
}

func (id AzureEndpointId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s/azureEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.TrafficManagerProfileName, id.Name)
}

// AzureEndpointID parses a AzureEndpoint ID into an AzureEndpointId struct
func AzureEndpointID(input string) (*AzureEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AzureEndpointId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.TrafficManagerProfileName, err = id.PopSegment("trafficManagerProfiles"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("azureEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
