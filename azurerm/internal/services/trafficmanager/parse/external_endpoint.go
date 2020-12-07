package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ExternalEndpointId struct {
	SubscriptionId            string
	ResourceGroup             string
	TrafficManagerProfileName string
	Name                      string
}

func NewExternalEndpointID(subscriptionId, resourceGroup, trafficManagerProfileName, name string) ExternalEndpointId {
	return ExternalEndpointId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		TrafficManagerProfileName: trafficManagerProfileName,
		Name:                      name,
	}
}

func (id ExternalEndpointId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s/externalEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.TrafficManagerProfileName, id.Name)
}

// ExternalEndpointID parses a ExternalEndpoint ID into an ExternalEndpointId struct
func ExternalEndpointID(input string) (*ExternalEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ExternalEndpointId{
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
	if resourceId.Name, err = id.PopSegment("externalEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
