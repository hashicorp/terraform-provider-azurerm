package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NestedEndpointId struct {
	SubscriptionId            string
	ResourceGroup             string
	TrafficManagerProfileName string
	Name                      string
}

func NewNestedEndpointID(subscriptionId, resourceGroup, trafficManagerProfileName, name string) NestedEndpointId {
	return NestedEndpointId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		TrafficManagerProfileName: trafficManagerProfileName,
		Name:                      name,
	}
}

func (id NestedEndpointId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Traffic Manager Profile Name %q", id.TrafficManagerProfileName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
}

func (id NestedEndpointId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s/nestedEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.TrafficManagerProfileName, id.Name)
}

// NestedEndpointID parses a NestedEndpoint ID into an NestedEndpointId struct
func NestedEndpointID(input string) (*NestedEndpointId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NestedEndpointId{
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
	if resourceId.Name, err = id.PopSegment("nestedEndpoints"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
