package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2020-03-13/services"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AdminKeysClient *adminkeys.AdminKeysClient
	QueryKeysClient *querykeys.QueryKeysClient
	ServicesClient  *services.ServicesClient
}

func NewClient(o *common.ClientOptions) *Client {
	adminKeysClient := adminkeys.NewAdminKeysClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&adminKeysClient.Client, o.ResourceManagerAuthorizer)

	queryKeysClient := querykeys.NewQueryKeysClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&queryKeysClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := services.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AdminKeysClient: &adminKeysClient,
		QueryKeysClient: &queryKeysClient,
		ServicesClient:  &servicesClient,
	}
}
