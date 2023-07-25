// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/caches"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-01-01/storagetargets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CachesClient         *caches.CachesClient
	StorageTargetsClient *storagetargets.StorageTargetsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	cachesClient, err := caches.NewCachesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Caches Client: %+v", err)
	}
	o.Configure(cachesClient.Client, o.Authorizers.ResourceManager)

	storageTargetsClient, err := storagetargets.NewStorageTargetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Storage Targets Client: %+v", err)
	}
	o.Configure(storageTargetsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CachesClient:         cachesClient,
		StorageTargetsClient: storageTargetsClient,
	}, nil
}
