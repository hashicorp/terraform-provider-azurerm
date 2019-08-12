package trafficmanager

import (
	"github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	GeographialHierarchiesClient trafficmanager.GeographicHierarchiesClient
	ProfilesClient               trafficmanager.ProfilesClient
	EndpointsClient              trafficmanager.EndpointsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.EndpointsClient = trafficmanager.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.EndpointsClient.Client, o.ResourceManagerAuthorizer)

	c.GeographialHierarchiesClient = trafficmanager.NewGeographicHierarchiesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GeographialHierarchiesClient.Client, o.ResourceManagerAuthorizer)

	c.ProfilesClient = trafficmanager.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProfilesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
