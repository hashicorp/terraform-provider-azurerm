package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CustomDomainsClient cdn.CustomDomainsClient
	EndpointsClient     cdn.EndpointsClient
	ProfilesClient      cdn.ProfilesClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.CustomDomainsClient = cdn.NewCustomDomainsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.CustomDomainsClient.Client, authorizer)

	c.EndpointsClient = cdn.NewEndpointsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.EndpointsClient.Client, authorizer)

	c.ProfilesClient = cdn.NewProfilesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ProfilesClient.Client, authorizer)

	return &c
}
