package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	CustomDomainsClient cdn.CustomDomainsClient
	EndpointsClient     cdn.EndpointsClient
	ProfilesClient      cdn.ProfilesClient
}

func BuildClients(endpoint, subscriptionId, partnerId string, auth autorest.Authorizer) *Client {
	c := Client{}

	c.CustomDomainsClient = cdn.NewCustomDomainsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.CustomDomainsClient.Client, auth, partnerId)

	c.EndpointsClient = cdn.NewEndpointsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.EndpointsClient.Client, auth, partnerId)

	c.ProfilesClient = cdn.NewProfilesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ProfilesClient.Client, auth, partnerId)

	return &c
}
