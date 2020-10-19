package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/storagecache/mgmt/2020-03-01/storagecache"
	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/tables"
)

type Client struct {
	AccountsClient           *storage.AccountsClient
	FileSystemsClient        *filesystems.Client
	ManagementPoliciesClient *storage.ManagementPoliciesClient
	BlobServicesClient       *storage.BlobServicesClient
	CachesClient             *storagecache.CachesClient
	StorageTargetsClient     *storagecache.StorageTargetsClient
	SyncServiceClient        *storagesync.ServicesClient
	SyncGroupsClient         *storagesync.SyncGroupsClient
	SubscriptionId           string

	environment   az.Environment
	storageAdAuth *autorest.Authorizer
}

func NewClient(options *common.ClientOptions) *Client {
	accountsClient := storage.NewAccountsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&accountsClient.Client, options.ResourceManagerAuthorizer)

	fileSystemsClient := filesystems.NewWithEnvironment(options.Environment)
	options.ConfigureClient(&fileSystemsClient.Client, options.StorageAuthorizer)

	managementPoliciesClient := storage.NewManagementPoliciesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&managementPoliciesClient.Client, options.ResourceManagerAuthorizer)

	blobServicesClient := storage.NewBlobServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&blobServicesClient.Client, options.ResourceManagerAuthorizer)

	cachesClient := storagecache.NewCachesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&cachesClient.Client, options.ResourceManagerAuthorizer)

	storageTargetsClient := storagecache.NewStorageTargetsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&storageTargetsClient.Client, options.ResourceManagerAuthorizer)

	syncServiceClient := storagesync.NewServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncServiceClient.Client, options.ResourceManagerAuthorizer)

	syncGroupsClient := storagesync.NewSyncGroupsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncGroupsClient.Client, options.ResourceManagerAuthorizer)

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	client := Client{
		AccountsClient:           &accountsClient,
		FileSystemsClient:        &fileSystemsClient,
		ManagementPoliciesClient: &managementPoliciesClient,
		BlobServicesClient:       &blobServicesClient,
		CachesClient:             &cachesClient,
		SubscriptionId:           options.SubscriptionId,
		StorageTargetsClient:     &storageTargetsClient,
		SyncServiceClient:        &syncServiceClient,
		SyncGroupsClient:         &syncGroupsClient,
		environment:              options.Environment,
	}

	if options.StorageUseAzureAD {
		client.storageAdAuth = &options.StorageAuthorizer
	}

	return &client
}

func (client Client) AccountsDataPlaneClient(ctx context.Context, account accountDetails) (*accounts.Client, error) {
	if client.storageAdAuth != nil {
		accountsClient := accounts.NewWithEnvironment(client.environment)
		accountsClient.Client.Authorizer = *client.storageAdAuth
		return &accountsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	accountsClient := accounts.NewWithEnvironment(client.environment)
	accountsClient.Client.Authorizer = storageAuth
	return &accountsClient, nil
}

func (client Client) BlobsClient(ctx context.Context, account accountDetails) (*blobs.Client, error) {
	if client.storageAdAuth != nil {
		blobsClient := blobs.NewWithEnvironment(client.environment)
		blobsClient.Client.Authorizer = *client.storageAdAuth
		return &blobsClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	blobsClient := blobs.NewWithEnvironment(client.environment)
	blobsClient.Client.Authorizer = storageAuth
	return &blobsClient, nil
}

func (client Client) ContainersClient(ctx context.Context, account accountDetails) (*containers.Client, error) {
	if client.storageAdAuth != nil {
		containersClient := containers.NewWithEnvironment(client.environment)
		containersClient.Client.Authorizer = *client.storageAdAuth
		return &containersClient, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	containersClient := containers.NewWithEnvironment(client.environment)
	containersClient.Client.Authorizer = storageAuth
	return &containersClient, nil
}

func (client Client) FileShareDirectoriesClient(ctx context.Context, account accountDetails) (*directories.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	directoriesClient := directories.NewWithEnvironment(client.environment)
	directoriesClient.Client.Authorizer = storageAuth
	return &directoriesClient, nil
}

func (client Client) FileSharesClient(ctx context.Context, account accountDetails) (*shares.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	sharesClient := shares.NewWithEnvironment(client.environment)
	sharesClient.Client.Authorizer = storageAuth
	return &sharesClient, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (*queues.Client, error) {
	if client.storageAdAuth != nil {
		queueAuth := queues.NewWithEnvironment(client.environment)
		queueAuth.Client.Authorizer = *client.storageAdAuth
		return &queueAuth, nil
	}

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	queuesClient := queues.NewWithEnvironment(client.environment)
	queuesClient.Client.Authorizer = storageAuth
	return &queuesClient, nil
}

func (client Client) TableEntityClient(ctx context.Context, account accountDetails) (*entities.Client, error) {
	// NOTE: Table Entity does not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	entitiesClient := entities.NewWithEnvironment(client.environment)
	entitiesClient.Client.Authorizer = storageAuth
	return &entitiesClient, nil
}

func (client Client) TablesClient(ctx context.Context, account accountDetails) (*tables.Client, error) {
	// NOTE: Tables do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	tablesClient := tables.NewWithEnvironment(client.environment)
	tablesClient.Client.Authorizer = storageAuth
	return &tablesClient, nil
}
