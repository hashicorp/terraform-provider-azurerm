package trafficmanager

import (
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GeographialHierarchiesClient *trafficmanager.GeographicHierarchiesClient
	ProfilesClient               *trafficmanager.ProfilesClient
	EndpointsClient              *trafficmanager.EndpointsClient
}

func BuildClient(o *common.ClientOptions) *Client {

	EndpointsClient := trafficmanager.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EndpointsClient.Client, o.ResourceManagerAuthorizer)

	GeographialHierarchiesClient := trafficmanager.NewGeographicHierarchiesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&GeographialHierarchiesClient.Client, o.ResourceManagerAuthorizer)

	ProfilesClient := trafficmanager.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProfilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		GeographialHierarchiesClient: &GeographialHierarchiesClient,
		ProfilesClient:               &ProfilesClient,
		EndpointsClient:              &EndpointsClient,
	}
}
