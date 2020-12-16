package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = EndpointId{}

// TODO: remove this shim layer once the resources have been split / in 3.0

type EndpointId struct {
	SubscriptionId            string
	ResourceGroup             string
	TrafficManagerProfileName string
	Name                      string
	endpointName              string
}

func (id EndpointId) EndpointType() string {
	return fmt.Sprintf("%sEndpoints", id.endpointName)
}

func NewEndpointId(subscriptionId, resourceGroup, trafficManagerProfile, endpoint, name string) (*EndpointId, error) {
	endpoint = strings.TrimSuffix(endpoint, "Endpoints")
	if endpoint != "azure" && endpoint != "external" && endpoint != "nested" {
		return nil, fmt.Errorf("unsupported value for Endpoint")
	}

	return &EndpointId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		TrafficManagerProfileName: trafficManagerProfile,
		Name:                      name,
		endpointName:              endpoint,
	}, nil
}

func (id EndpointId) ID() string {
	if id.endpointName == "azure" {
		return NewAzureEndpointID(id.SubscriptionId, id.ResourceGroup, id.TrafficManagerProfileName, id.Name).ID()
	}

	if id.endpointName == "external" {
		return NewExternalEndpointID(id.SubscriptionId, id.ResourceGroup, id.TrafficManagerProfileName, id.Name).ID()
	}

	if id.endpointName == "nested" {
		return NewNestedEndpointID(id.SubscriptionId, id.ResourceGroup, id.TrafficManagerProfileName, id.Name).ID()
	}

	// NOTE: this looks bad but is caught above in the `New` function and below in the Parse
	panic("not implemented")
}

func EndpointID(input string) (*EndpointId, error) {
	azureEndpoint, err := AzureEndpointID(input)
	if err == nil {
		return &EndpointId{
			SubscriptionId:            azureEndpoint.SubscriptionId,
			ResourceGroup:             azureEndpoint.ResourceGroup,
			TrafficManagerProfileName: azureEndpoint.TrafficManagerProfileName,
			Name:                      azureEndpoint.Name,
			endpointName:              "azure",
		}, nil
	}

	externalEndpoint, err := ExternalEndpointID(input)
	if err == nil {
		return &EndpointId{
			SubscriptionId:            externalEndpoint.SubscriptionId,
			ResourceGroup:             externalEndpoint.ResourceGroup,
			TrafficManagerProfileName: externalEndpoint.TrafficManagerProfileName,
			Name:                      externalEndpoint.Name,
			endpointName:              "external",
		}, nil
	}

	nestedEndpoint, err := NestedEndpointID(input)
	if err != nil {
		return nil, err
	}
	return &EndpointId{
		SubscriptionId:            nestedEndpoint.SubscriptionId,
		ResourceGroup:             nestedEndpoint.ResourceGroup,
		TrafficManagerProfileName: nestedEndpoint.TrafficManagerProfileName,
		Name:                      nestedEndpoint.Name,
		endpointName:              "nested",
	}, nil
}
