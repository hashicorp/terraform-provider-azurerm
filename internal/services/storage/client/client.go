// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	storage_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/paths"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/files"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/table/tables"
)

type Client struct {
	AccountsClient              *storage.AccountsClient
	FileSystemsClient           *filesystems.Client
	ADLSGen2PathsClient         *paths.Client
	BlobServicesClient          *storage.BlobServicesClient
	BlobInventoryPoliciesClient *storage.BlobInventoryPoliciesClient
	EncryptionScopesClient      *storage.EncryptionScopesClient
	Environment                 azure.Environment
	FileServicesClient          *storage.FileServicesClient
	SyncCloudEndpointsClient    *cloudendpointresource.CloudEndpointResourceClient
	SyncServiceClient           *storagesyncservicesresource.StorageSyncServicesResourceClient
	SyncGroupsClient            *syncgroupresource.SyncGroupResourceClient
	SubscriptionId              string

	ResourceManager *storage_v2023_01_01.Client

	resourceManagerAuthorizer autorest.Authorizer
	storageAdAuth             *autorest.Authorizer
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	accountsClient := storage.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountsClient.Client, o.ResourceManagerAuthorizer)

	fileSystemsClient := filesystems.NewWithEnvironment(o.AzureEnvironment)
	o.ConfigureClient(&fileSystemsClient.Client, o.StorageAuthorizer)

	adlsGen2PathsClient := paths.NewWithEnvironment(o.AzureEnvironment)
	o.ConfigureClient(&adlsGen2PathsClient.Client, o.StorageAuthorizer)

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
		FileSystemsClient:           &fileSystemsClient,
		ADLSGen2PathsClient:         &adlsGen2PathsClient,
		BlobServicesClient:          &blobServicesClient,
		BlobInventoryPoliciesClient: &blobInventoryPoliciesClient,
		EncryptionScopesClient:      &encryptionScopesClient,
		Environment:                 o.AzureEnvironment,
		FileServicesClient:          &fileServicesClient,
		ResourceManager:             resourceManager,
		SubscriptionId:              o.SubscriptionId,
		SyncCloudEndpointsClient:    syncCloudEndpointsClient,
		SyncServiceClient:           syncServiceClient,
		SyncGroupsClient:            syncGroupsClient,

		resourceManagerAuthorizer: o.ResourceManagerAuthorizer,
	}

	if o.StorageUseAzureAD {
		client.storageAdAuth = &o.StorageAuthorizer
	}

	return &client, nil
}

func (client Client) AccountsDataPlaneClient(ctx context.Context, account accountDetails) (*accounts.Client, error) {
	if client.storageAdAuth != nil {
		accountsClient := accounts.NewWithEnvironment(client.Environment)
		accountsClient.Client.Authorizer = *client.storageAdAuth
		return &accountsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	accountsClient := accounts.NewWithEnvironment(client.Environment)
	accountsClient.Client.Authorizer = storageAuth
	return &accountsClient, nil
}

func (client Client) BlobsClient(ctx context.Context, account accountDetails) (*blobs.Client, error) {
	if client.storageAdAuth != nil {
		blobsClient := blobs.NewWithEnvironment(client.Environment)
		blobsClient.Client.Authorizer = *client.storageAdAuth
		return &blobsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	blobsClient := blobs.NewWithEnvironment(client.Environment)
	blobsClient.Client.Authorizer = storageAuth
	return &blobsClient, nil
}

func (client Client) ContainersClient(ctx context.Context, account accountDetails) (shim.StorageContainerWrapper, error) {
	if client.storageAdAuth != nil {
		containersClient := containers.NewWithEnvironment(client.Environment)
		containersClient.Client.Authorizer = *client.storageAdAuth
		shim := shim.NewDataPlaneStorageContainerWrapper(&containersClient)
		return shim, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	containersClient := containers.NewWithEnvironment(client.Environment)
	containersClient.Client.Authorizer = storageAuth

	shim := shim.NewDataPlaneStorageContainerWrapper(&containersClient)
	return shim, nil
}

func (client Client) FileShareDirectoriesClient(ctx context.Context, account accountDetails) (*directories.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	directoriesClient := directories.NewWithEnvironment(client.Environment)
	directoriesClient.Client.Authorizer = storageAuth
	return &directoriesClient, nil
}

func (client Client) FileShareFilesClient(ctx context.Context, account accountDetails) (*files.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	filesClient := files.NewWithEnvironment(client.Environment)
	filesClient.Client.Authorizer = storageAuth
	return &filesClient, nil
}

func (client Client) FileSharesClient(ctx context.Context, account accountDetails) (shim.StorageShareWrapper, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	sharesClient := shares.NewWithEnvironment(client.Environment)
	sharesClient.Client.Authorizer = storageAuth
	shim := shim.NewDataPlaneStorageShareWrapper(&sharesClient)
	return shim, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (shim.StorageQueuesWrapper, error) {
	if client.storageAdAuth != nil {
		queueClient := queues.NewWithEnvironment(client.Environment)
		queueClient.Client.Authorizer = *client.storageAdAuth
		return shim.NewDataPlaneStorageQueueWrapper(&queueClient), nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	queuesClient := queues.NewWithEnvironment(client.Environment)
	queuesClient.Client.Authorizer = storageAuth
	return shim.NewDataPlaneStorageQueueWrapper(&queuesClient), nil
}

func (client Client) TableEntityClient(ctx context.Context, account accountDetails) (*entities.Client, error) {
	// NOTE: Table Entity does not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	entitiesClient := entities.NewWithEnvironment(client.Environment)
	entitiesClient.Client.Authorizer = storageAuth
	return &entitiesClient, nil
}

func (client Client) TablesClient(ctx context.Context, account accountDetails) (shim.StorageTableWrapper, error) {
	// NOTE: Tables do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	tablesClient := tables.NewWithEnvironment(client.Environment)
	tablesClient.Client.Authorizer = storageAuth
	shim := shim.NewDataPlaneStorageTableWrapper(&tablesClient)
	return shim, nil
}
