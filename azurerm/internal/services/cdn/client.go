package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CustomDomainsClient cdn.CustomDomainsClient
	EndpointsClient     cdn.EndpointsClient
	ProfilesClient      cdn.ProfilesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.CustomDomainsClient = cdn.NewCustomDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CustomDomainsClient.Client, o.ResourceManagerAuthorizer)

	c.EndpointsClient = cdn.NewEndpointsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.EndpointsClient.Client, o.ResourceManagerAuthorizer)

	c.ProfilesClient = cdn.NewProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProfilesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
