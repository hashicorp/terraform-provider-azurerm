package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CustomDomainsClient *cdn.CustomDomainsClient
	EndpointsClient     *cdn.EndpointsClient
	ProfilesClient      *cdn.ProfilesClient
}

func NewClient(o *common.ClientOptions) *Client {
	customDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&customDomainsClient.Client, o.ResourceManagerAuthorizer)

	endpointsClient := cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&endpointsClient.Client, o.ResourceManagerAuthorizer)

	profilesClient := cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&profilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CustomDomainsClient: &customDomainsClient,
		EndpointsClient:     &endpointsClient,
		ProfilesClient:      &profilesClient,
	}
}
