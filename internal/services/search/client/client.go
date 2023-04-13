package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/adminkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/querykeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/services"
	"github.com/hashicorp/go-azure-sdk/resource-manager/search/2022-09-01/sharedprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AdminKeysClient                       *adminkeys.AdminKeysClient
	QueryKeysClient                       *querykeys.QueryKeysClient
	ServicesClient                        *services.ServicesClient
	SearchSharedPrivateLinkResourceClient *sharedprivatelinkresources.SharedPrivateLinkResourcesClient
}

func NewClient(o *common.ClientOptions) *Client {
	adminKeysClient := adminkeys.NewAdminKeysClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&adminKeysClient.Client, o.ResourceManagerAuthorizer)

	queryKeysClient := querykeys.NewQueryKeysClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&queryKeysClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := services.NewServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	searchSharedPrivateLinkResourceClient := sharedprivatelinkresources.NewSharedPrivateLinkResourcesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&searchSharedPrivateLinkResourceClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AdminKeysClient:                       &adminKeysClient,
		QueryKeysClient:                       &queryKeysClient,
		ServicesClient:                        &servicesClient,
		SearchSharedPrivateLinkResourceClient: &searchSharedPrivateLinkResourceClient,
	}
}
