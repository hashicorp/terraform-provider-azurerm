package search

import (
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AdminKeysClient search.AdminKeysClient
	ServicesClient  search.ServicesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AdminKeysClient = search.NewAdminKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AdminKeysClient.Client, o.ResourceManagerAuthorizer)

	c.ServicesClient = search.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServicesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
