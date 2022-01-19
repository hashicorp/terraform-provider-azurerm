package client

import (
	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2021-09-01/storagecache"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CachesClient         *storagecache.CachesClient
	StorageTargetsClient *storagecache.StorageTargetsClient
}

func NewClient(options *common.ClientOptions) *Client {
	cachesClient := storagecache.NewCachesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&cachesClient.Client, options.ResourceManagerAuthorizer)

	storageTargetsClient := storagecache.NewStorageTargetsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&storageTargetsClient.Client, options.ResourceManagerAuthorizer)

	return &Client{
		CachesClient:         &cachesClient,
		StorageTargetsClient: &storageTargetsClient,
	}
}
