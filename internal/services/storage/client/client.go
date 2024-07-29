// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	storage_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/registeredserverresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

// StorageDomainSuffix is used by validation functions
var StorageDomainSuffix *string

type Client struct {
	StorageDomainSuffix string

	ResourceManager *storage_v2023_01_01.Client
	// TODO: import the Storage Sync Meta Client and use that
	SyncCloudEndpointsClient   *cloudendpointresource.CloudEndpointResourceClient
	SyncGroupsClient           *syncgroupresource.SyncGroupResourceClient
	SyncRegisteredServerClient *registeredserverresource.RegisteredServerResourceClient
	SyncServerEndpointsClient  *serverendpointresource.ServerEndpointResourceClient
	SyncServiceClient          *storagesyncservicesresource.StorageSyncServicesResourceClient

	authConfigForAzureAD *auth.Credentials
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	storageSuffix, ok := o.Environment.Storage.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("determining domain suffix for storage in environment: %s", o.Environment.Name)
	}

	// Set global variable for post-configure validation
	StorageDomainSuffix = storageSuffix

	resourceManager, err := storage_v2023_01_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building ResourceManager clients: %+v", err)
	}

	syncCloudEndpointsClient, err := cloudendpointresource.NewCloudEndpointResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building CloudEndpoint client: %+v", err)
	}
	o.Configure(syncCloudEndpointsClient.Client, o.Authorizers.ResourceManager)

	syncRegisteredServersClient, err := registeredserverresource.NewRegisteredServerResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageRegisteredServer client: %+v", err)
	}
	o.Configure(syncRegisteredServersClient.Client, o.Authorizers.ResourceManager)

	syncServerEndpointClient, err := serverendpointresource.NewServerEndpointResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageSyncServerEndpoint client: %+v", err)
	}
	o.Configure(syncServerEndpointClient.Client, o.Authorizers.ResourceManager)

	syncServiceClient, err := storagesyncservicesresource.NewStorageSyncServicesResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageSyncService client: %+v", err)
	}
	o.Configure(syncServiceClient.Client, o.Authorizers.ResourceManager)

	syncGroupsClient, err := syncgroupresource.NewSyncGroupResourceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building StorageSyncGroups client: %+v", err)
	}
	o.Configure(syncGroupsClient.Client, o.Authorizers.ResourceManager)

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	client := Client{
		ResourceManager:            resourceManager,
		SyncCloudEndpointsClient:   syncCloudEndpointsClient,
		SyncRegisteredServerClient: syncRegisteredServersClient,
		SyncServerEndpointsClient:  syncServerEndpointClient,
		SyncServiceClient:          syncServiceClient,
		SyncGroupsClient:           syncGroupsClient,

		StorageDomainSuffix: *storageSuffix,
	}

	if o.StorageUseAzureAD {
		client.authConfigForAzureAD = o.AuthConfig
	}

	return &client, nil
}
