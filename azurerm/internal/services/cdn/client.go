package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	CustomDomainsClient cdn.CustomDomainsClient
	EndpointsClient     cdn.EndpointsClient
	ProfilesClient      cdn.ProfilesClient
}

func BuildClient(endpoint, subscriptionId string, o *ar.ClientOptions) *Client {
	c := Client{}

	c.CustomDomainsClient = cdn.NewCustomDomainsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.CustomDomainsClient.Client, o)

	c.EndpointsClient = cdn.NewEndpointsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.EndpointsClient.Client, o)

	c.ProfilesClient = cdn.NewProfilesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProfilesClient.Client, o)

	return &c
}
