package client

import (
	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2015-08-19/search"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AdminKeysClient *search.AdminKeysClient
	ServicesClient  *search.ServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	adminKeysClient := search.NewAdminKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&adminKeysClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := search.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AdminKeysClient: &adminKeysClient,
		ServicesClient:  &servicesClient,
	}
}
