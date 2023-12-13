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
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
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
	storageAdAuth             *auth.Authorizer
}

func NewClient(o *common.ClientOptions) (*Client, error) {
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
		accountsClient, err := accounts.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Blob)
		if err != nil {
			return nil, fmt.Errorf("creating Accounts Client: %+v", err)
		}
		accountsClient.Client.WithAuthorizer(*client.storageAdAuth)
		return accountsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	accountsClient, err := accounts.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Blob)
	if err != nil {
		return nil, fmt.Errorf("creating Accounts Client: %+v", err)
	}
	accountsClient.Client.WithAuthorizer(storageAuth)
	return accountsClient, nil
}

func (client Client) BlobsClient(ctx context.Context, account accountDetails) (*blobs.Client, error) {
	if client.storageAdAuth != nil {
		blobsClient, err := blobs.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Blob)
		if err != nil {
			return nil, fmt.Errorf("creating Blobs Client: %+v", err)
		}
		blobsClient.Client.WithAuthorizer(*client.storageAdAuth)
		return blobsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	blobsClient, err := blobs.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Blob)
	if err != nil {
		return nil, fmt.Errorf("creating Blobs Client: %+v", err)
	}
	blobsClient.Client.WithAuthorizer(storageAuth)
	return blobsClient, nil
}

func (client Client) ContainersClient(ctx context.Context, account accountDetails) (shim.StorageContainerWrapper, error) {
	if client.storageAdAuth != nil {
		containersClient, err := containers.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Blob)
		if err != nil {
			return nil, fmt.Errorf("creating Containers Client: %+v", err)
		}
		containersClient.Client.WithAuthorizer(*client.storageAdAuth)
		shim := shim.NewDataPlaneStorageContainerWrapper(containersClient)
		return shim, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	containersClient, err := containers.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Blob)
	if err != nil {
		return nil, fmt.Errorf("creating Containers Client: %+v", err)
	}
	containersClient.Client.WithAuthorizer(storageAuth)

	shim := shim.NewDataPlaneStorageContainerWrapper(containersClient)
	return shim, nil
}

func (client Client) FileShareDirectoriesClient(ctx context.Context, account accountDetails) (*directories.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	directoriesClient, err := directories.NewWithBaseUri(*account.Properties.PrimaryEndpoints.File)
	if err != nil {
		return nil, fmt.Errorf("creating Directories Client: %+v", err)
	}
	directoriesClient.Client.WithAuthorizer(storageAuth)
	return directoriesClient, nil
}

func (client Client) FileShareFilesClient(ctx context.Context, account accountDetails) (*files.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	filesClient, err := files.NewWithBaseUri(*account.Properties.PrimaryEndpoints.File)
	if err != nil {
		return nil, fmt.Errorf("creating Files Client: %+v", err)
	}
	filesClient.Client.WithAuthorizer(storageAuth)
	return filesClient, nil
}

func (client Client) FileSharesClient(ctx context.Context, account accountDetails) (shim.StorageShareWrapper, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	sharesClient, err := shares.NewWithBaseUri(*account.Properties.PrimaryEndpoints.File)
	if err != nil {
		return nil, fmt.Errorf("creating Shares Client: %+v", err)
	}
	sharesClient.Client.WithAuthorizer(storageAuth)
	shim := shim.NewDataPlaneStorageShareWrapper(sharesClient)
	return shim, nil
}

func (client Client) FileSystemsClient(ctx context.Context, account accountDetails) (*filesystems.Client, error) {
	if client.storageAdAuth != nil {
		filesystemsClient, err := filesystems.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Dfs)
		if err != nil {
			return nil, fmt.Errorf("creating Filesystems Client: %+v", err)
		}
		filesystemsClient.Client.WithAuthorizer(*client.storageAdAuth)
		return filesystemsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	filesystemsClient, err := filesystems.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Dfs)
	if err != nil {
		return nil, fmt.Errorf("creating Filesystems Client: %+v", err)
	}
	filesystemsClient.Client.WithAuthorizer(storageAuth)
	return filesystemsClient, nil
}

func (client Client) PathsClient(ctx context.Context, account accountDetails) (*paths.Client, error) {
	if client.storageAdAuth != nil {
		pathsClient, err := paths.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Dfs)
		if err != nil {
			return nil, fmt.Errorf("creating Filesystems Client: %+v", err)
		}
		pathsClient.Client.WithAuthorizer(*client.storageAdAuth)
		return pathsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	pathsClient, err := paths.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Dfs)
	if err != nil {
		return nil, fmt.Errorf("creating Filesystems Client: %+v", err)
	}
	pathsClient.Client.WithAuthorizer(storageAuth)
	return pathsClient, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (shim.StorageQueuesWrapper, error) {
	if client.storageAdAuth != nil {
		queueClient, err := queues.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Queue)
		if err != nil {
			return nil, fmt.Errorf("creating Queues Client: %+v", err)
		}
		queueClient.Client.WithAuthorizer(*client.storageAdAuth)
		return shim.NewDataPlaneStorageQueueWrapper(queueClient), nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	queuesClient, err := queues.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Queue)
	if err != nil {
		return nil, fmt.Errorf("creating Queues Client: %+v", err)
	}
	queuesClient.Client.WithAuthorizer(storageAuth)
	return shim.NewDataPlaneStorageQueueWrapper(queuesClient), nil
}

func (client Client) TableEntityClient(ctx context.Context, account accountDetails) (*entities.Client, error) {
	// NOTE: Table Entity does not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKeyTable)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	entitiesClient, err := entities.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Table)
	if err != nil {
		return nil, fmt.Errorf("creating Queues Client: %+v", err)
	}
	entitiesClient.Client.WithAuthorizer(storageAuth)
	return entitiesClient, nil
}

func (client Client) TablesClient(ctx context.Context, account accountDetails) (shim.StorageTableWrapper, error) {
	// NOTE: Tables do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("retrieving Account Key: %s", err)
	}

	storageAuth, err := auth.NewSharedKeyAuthorizer(account.name, *accountKey, auth.SharedKeyTable)
	if err != nil {
		return nil, fmt.Errorf("building Authorizer: %+v", err)
	}

	tablesClient, err := tables.NewWithBaseUri(*account.Properties.PrimaryEndpoints.Table)
	if err != nil {
		return nil, fmt.Errorf("creating Queues Client: %+v", err)
	}
	tablesClient.Client.WithAuthorizer(storageAuth)
	shim := shim.NewDataPlaneStorageTableWrapper(tablesClient)
	return shim, nil
}
