package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CustomDomainsClient *cdn.CustomDomainsClient
	EndpointsClient     *cdn.EndpointsClient
	ProfilesClient      *cdn.ProfilesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	CustomDomainsClient := cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&CustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	EndpointsClient := cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EndpointsClient.Client, o.ResourceManagerAuthorizer)

	ProfilesClient := cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ProfilesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CustomDomainsClient: &CustomDomainsClient,
		EndpointsClient:     &EndpointsClient,
		ProfilesClient:      &ProfilesClient,
	}
}
