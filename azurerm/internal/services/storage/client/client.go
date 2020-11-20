package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/shim"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/paths"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/files"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

type Client struct {
	AccountsClient           *storage.AccountsClient
	FileSystemsClient        *filesystems.Client
	ADLSGen2PathsClient      *paths.Client
	ManagementPoliciesClient *storage.ManagementPoliciesClient
	BlobServicesClient       *storage.BlobServicesClient
	EncryptionScopesClient   *storage.EncryptionScopesClient
	Environment              az.Environment
	SyncServiceClient        *storagesync.ServicesClient
	SyncGroupsClient         *storagesync.SyncGroupsClient
	SubscriptionId           string

	resourceManagerAuthorizer autorest.Authorizer
	storageAdAuth             *autorest.Authorizer
}

func NewClient(options *common.ClientOptions) *Client {
	accountsClient := storage.NewAccountsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&accountsClient.Client, options.ResourceManagerAuthorizer)

	fileSystemsClient := filesystems.NewWithEnvironment(options.Environment)
	options.ConfigureClient(&fileSystemsClient.Client, options.StorageAuthorizer)

	adlsGen2PathsClient := paths.NewWithEnvironment(options.Environment)
	options.ConfigureClient(&adlsGen2PathsClient.Client, options.StorageAuthorizer)

	managementPoliciesClient := storage.NewManagementPoliciesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&managementPoliciesClient.Client, options.ResourceManagerAuthorizer)

	blobServicesClient := storage.NewBlobServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&blobServicesClient.Client, options.ResourceManagerAuthorizer)

	encryptionScopesClient := storage.NewEncryptionScopesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&encryptionScopesClient.Client, options.ResourceManagerAuthorizer)

	syncServiceClient := storagesync.NewServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncServiceClient.Client, options.ResourceManagerAuthorizer)

	syncGroupsClient := storagesync.NewSyncGroupsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncGroupsClient.Client, options.ResourceManagerAuthorizer)

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	client := Client{
		AccountsClient:           &accountsClient,
		FileSystemsClient:        &fileSystemsClient,
		ADLSGen2PathsClient:      &adlsGen2PathsClient,
		ManagementPoliciesClient: &managementPoliciesClient,
		BlobServicesClient:       &blobServicesClient,
		EncryptionScopesClient:   &encryptionScopesClient,
		Environment:              options.Environment,
		SubscriptionId:           options.SubscriptionId,
		SyncServiceClient:        &syncServiceClient,
		SyncGroupsClient:         &syncGroupsClient,

		resourceManagerAuthorizer: options.ResourceManagerAuthorizer,
	}

	if options.StorageUseAzureAD {
		client.storageAdAuth = &options.StorageAuthorizer
	}

	return &client
}

func (client Client) AccountsDataPlaneClient(ctx context.Context, account accountDetails) (*accounts.Client, error) {
	if client.storageAdAuth != nil {
		accountsClient := accounts.NewWithEnvironment(client.Environment)
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
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
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
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKey)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
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
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	directoriesClient := directories.NewWithEnvironment(client.Environment)
	directoriesClient.Client.Authorizer = storageAuth
	return &directoriesClient, nil
}

func (client Client) FileShareFilesClient(ctx context.Context, account accountDetails) (*files.Client, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	filesClient := files.NewWithEnvironment(client.Environment)
	filesClient.Client.Authorizer = storageAuth
	return &filesClient, nil
}

func (client Client) FileSharesClient(ctx context.Context, account accountDetails) (shim.StorageShareWrapper, error) {
	// NOTE: Files do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
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
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	entitiesClient := entities.NewWithEnvironment(client.Environment)
	entitiesClient.Client.Authorizer = storageAuth
	return &entitiesClient, nil
}

func (client Client) TablesClient(ctx context.Context, account accountDetails) (shim.StorageTableWrapper, error) {
	// NOTE: Tables do not support AzureAD Authentication

	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLiteForTable)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	tablesClient := tables.NewWithEnvironment(client.Environment)
	tablesClient.Client.Authorizer = storageAuth
	shim := shim.NewDataPlaneStorageTableWrapper(&tablesClient)
	return shim, nil
}
