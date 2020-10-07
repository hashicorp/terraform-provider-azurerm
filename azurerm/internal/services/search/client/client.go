package client

import (
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2020-03-13/search"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AdminKeysClient *search.AdminKeysClient
	QueryKeysClient *search.QueryKeysClient
	ServicesClient  *search.ServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	adminKeysClient := search.NewAdminKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&adminKeysClient.Client, o.ResourceManagerAuthorizer)

	queryKeysClient := search.NewQueryKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&queryKeysClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := search.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AdminKeysClient: &adminKeysClient,
		QueryKeysClient: &queryKeysClient,
		ServicesClient:  &servicesClient,
	}
}
