package client

import (
	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2019-11-01/storagecache"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CachesClient *storagecache.CachesClient
}

func NewClient(o *common.ClientOptions) *Client {
	cachesClient := storagecache.NewCachesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cachesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CachesClient: &cachesClient,
	}
}
