package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/resourcemanagershim"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
)

type Client struct {
	AccountsClient           *storage.AccountsClient
	FileSystemsClient        *filesystems.Client
	ManagementPoliciesClient *storage.ManagementPoliciesClient
	BlobServicesClient       *storage.BlobServicesClient
	Environment              az.Environment
	SyncServiceClient        *storagesync.ServicesClient
	SyncGroupsClient         *storagesync.SyncGroupsClient
	SubscriptionId           string

	resourceManagerAuthorizer autorest.Authorizer
	storageAdAuth             *autorest.Authorizer
	useResourceManager        bool
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
		Environment:              options.Environment,
		SubscriptionId:           options.SubscriptionId,
		SyncServiceClient:        &syncServiceClient,
		SyncGroupsClient:         &syncGroupsClient,

		resourceManagerAuthorizer: options.ResourceManagerAuthorizer,

		useResourceManager: false, // TODO: feature toggleable
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

func (client Client) ContainersClient(ctx context.Context, account accountDetails) (resourcemanagershim.StorageContainerWrapper, error) {
	if client.useResourceManager {
		rmClient := storage.NewBlobContainersClientWithBaseURI(client.Environment.ResourceManagerEndpoint, client.SubscriptionId)
		rmClient.Client.Authorizer = client.resourceManagerAuthorizer
		rmShim := resourcemanagershim.NewResourceManagerStorageContainerWrapper(&rmClient)
		return &rmShim, nil
	}

	if client.storageAdAuth != nil {
		containersClient := containers.NewWithEnvironment(client.Environment)
		containersClient.Client.Authorizer = *client.storageAdAuth
		shim := resourcemanagershim.NewDataPlaneStorageContainerWrapper(&containersClient)
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

	shim := resourcemanagershim.NewDataPlaneStorageContainerWrapper(&containersClient)
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

	sharesClient := shares.NewWithEnvironment(client.Environment)
	sharesClient.Client.Authorizer = storageAuth
	return &sharesClient, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (*queues.Client, error) {
	if client.storageAdAuth != nil {
		queueAuth := queues.NewWithEnvironment(client.Environment)
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

	queuesClient := queues.NewWithEnvironment(client.Environment)
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

	entitiesClient := entities.NewWithEnvironment(client.Environment)
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

	tablesClient := tables.NewWithEnvironment(client.Environment)
	tablesClient.Client.Authorizer = storageAuth
	return &tablesClient, nil
}
