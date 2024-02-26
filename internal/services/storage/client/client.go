// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/azure"
	storage_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SubscriptionId string

	AccountsClient              *storage.AccountsClient
	BlobInventoryPoliciesClient *storage.BlobInventoryPoliciesClient
	BlobServicesClient          *storage.BlobServicesClient
	EncryptionScopesClient      *storage.EncryptionScopesClient
	FileServicesClient          *storage.FileServicesClient
	SyncCloudEndpointsClient    *cloudendpointresource.CloudEndpointResourceClient
	SyncGroupsClient            *syncgroupresource.SyncGroupResourceClient
	SyncServiceClient           *storagesyncservicesresource.StorageSyncServicesResourceClient

	ResourceManager *storage_v2023_01_01.Client

	AzureEnvironment    azure.Environment
	StorageDomainSuffix string

	authorizerForAad auth.Authorizer
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	storageSuffix, ok := o.Environment.Storage.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("determining domain suffix for storage in environment: %s", o.Environment.Name)
	}

	accountsClient := storage.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	blobServicesClient := storage.NewBlobServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&blobServicesClient.Client, o.ResourceManagerAuthorizer)

	blobInventoryPoliciesClient := storage.NewBlobInventoryPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&blobInventoryPoliciesClient.Client, o.ResourceManagerAuthorizer)

	encryptionScopesClient := storage.NewEncryptionScopesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&encryptionScopesClient.Client, o.ResourceManagerAuthorizer)

	fileServicesClient := storage.NewFileServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&fileServicesClient.Client, o.ResourceManagerAuthorizer)

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
		AccountsClient:              &accountsClient,
		BlobServicesClient:          &blobServicesClient,
		BlobInventoryPoliciesClient: &blobInventoryPoliciesClient,
		EncryptionScopesClient:      &encryptionScopesClient,
		FileServicesClient:          &fileServicesClient,
		ResourceManager:             resourceManager,
		SubscriptionId:              o.SubscriptionId,
		SyncCloudEndpointsClient:    syncCloudEndpointsClient,
		SyncServiceClient:           syncServiceClient,
		SyncGroupsClient:            syncGroupsClient,

		AzureEnvironment:    o.AzureEnvironment,
		StorageDomainSuffix: *storageSuffix,
	}

	if o.StorageUseAzureAD {
		client.authorizerForAad = o.Authorizers.Storage
	}

	return &client, nil
}
