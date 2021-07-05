package client

import (
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-08-01/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GeographialHierarchiesClient *trafficmanager.GeographicHierarchiesClient
	EndpointsClient              *trafficmanager.EndpointsClient
	ProfilesClient               *trafficmanager.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	endpointsClient := trafficmanager.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	geographialHierarchiesClient := trafficmanager.NewGeographicHierarchiesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&geographialHierarchiesClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := trafficmanager.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		EndpointsClient:              &endpointsClient,
		GeographialHierarchiesClient: &geographialHierarchiesClient,
		ProfilesClient:               &profilesClient,
	}
}
