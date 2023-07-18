// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/caches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/storagetargets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01/amlfilesystems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AMLFileSystemsClient *amlfilesystems.AmlFilesystemsClient
	CachesClient         *caches.CachesClient
	StorageTargetsClient *storagetargets.StorageTargetsClient
}

func NewClient(options *common.ClientOptions) *Client {
	amlFileSystemsClient := amlfilesystems.NewAmlFilesystemsClientWithBaseURI(options.ResourceManagerEndpoint)
	options.ConfigureClient(&amlFileSystemsClient.Client, options.ResourceManagerAuthorizer)

	cachesClient := caches.NewCachesClientWithBaseURI(options.ResourceManagerEndpoint)
	options.ConfigureClient(&cachesClient.Client, options.ResourceManagerAuthorizer)

	storageTargetsClient := storagetargets.NewStorageTargetsClientWithBaseURI(options.ResourceManagerEndpoint)
	options.ConfigureClient(&storageTargetsClient.Client, options.ResourceManagerAuthorizer)

	return &Client{
		AMLFileSystemsClient: &amlFileSystemsClient,
		CachesClient:         &cachesClient,
		StorageTargetsClient: &storageTargetsClient,
	}
}
