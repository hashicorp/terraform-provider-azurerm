package search

import (
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AdminKeysClient *search.AdminKeysClient
	ServicesClient  *search.ServicesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	AdminKeysClient := search.NewAdminKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AdminKeysClient.Client, o.ResourceManagerAuthorizer)

	ServicesClient := search.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AdminKeysClient: &AdminKeysClient,
		ServicesClient:  &ServicesClient,
	}
}
