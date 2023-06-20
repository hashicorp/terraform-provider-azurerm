package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/caches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CachesClient         *caches.CachesClient
	StorageTargetsClient *storagetargets.StorageTargetsClient
}

func NewClient(options *common.ClientOptions) *Client {

	cachesClient := caches.NewCachesClientWithBaseURI(options.ResourceManagerEndpoint)
	options.ConfigureClient(&cachesClient.Client, options.ResourceManagerAuthorizer)

	storageTargetsClient := storagetargets.NewStorageTargetsClientWithBaseURI(options.ResourceManagerEndpoint)
	options.ConfigureClient(&storageTargetsClient.Client, options.ResourceManagerAuthorizer)

	return &Client{
		CachesClient:         &cachesClient,
		StorageTargetsClient: &storageTargetsClient,
	}
}
