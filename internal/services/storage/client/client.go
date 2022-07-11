package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync"
	"github.com/Azure/go-autorest/autorest"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/sdk/2021-04-01/objectreplicationpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/shim"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/paths"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/table/tables"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/files"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
)

type Client struct {
	AccountsClient              *storage.AccountsClient
	FileSystemsClient           *filesystems.Client
	ADLSGen2PathsClient         *paths.Client
	ManagementPoliciesClient    *storage.ManagementPoliciesClient
	BlobInventoryPoliciesClient *storage.BlobInventoryPoliciesClient
	CloudEndpointsClient        *storagesync.CloudEndpointsClient
	EncryptionScopesClient      *storage.EncryptionScopesClient
	Environment                 az.Environment
	ObjectReplicationClient     *objectreplicationpolicies.ObjectReplicationPoliciesClient
	SyncServiceClient           *storagesync.ServicesClient
	SyncGroupsClient            *storagesync.SyncGroupsClient
	SubscriptionId              string

	BlobServicesClient  *storage.BlobServicesClient
	FileServicesClient  *storage.FileServicesClient
	QueueServicesClient *storage.QueueServicesClient
	TableServicesClient *storage.TableServicesClient

	resourceManagerAuthorizer autorest.Authorizer
	storageAdAuth             *autorest.Authorizer

	// useResourceManager specifies whether to use the mgmt plane API for resources that have both data plane and mgmt plane support.
	// Currently, only the following resources are affected: blob container, file share, queue.
	// TODO: table.
	useResourceManager bool
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

	blobInventoryPoliciesClient := storage.NewBlobInventoryPoliciesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&blobInventoryPoliciesClient.Client, options.ResourceManagerAuthorizer)

	cloudEndpointsClient := storagesync.NewCloudEndpointsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&cloudEndpointsClient.Client, options.ResourceManagerAuthorizer)

	encryptionScopesClient := storage.NewEncryptionScopesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&encryptionScopesClient.Client, options.ResourceManagerAuthorizer)

	objectReplicationPolicyClient := objectreplicationpolicies.NewObjectReplicationPoliciesClientWithBaseURI(options.ResourceManagerEndpoint)
	options.ConfigureClient(&objectReplicationPolicyClient.Client, options.ResourceManagerAuthorizer)

	syncServiceClient := storagesync.NewServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncServiceClient.Client, options.ResourceManagerAuthorizer)

	syncGroupsClient := storagesync.NewSyncGroupsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncGroupsClient.Client, options.ResourceManagerAuthorizer)

	blobServicesClient := storage.NewBlobServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&blobServicesClient.Client, options.ResourceManagerAuthorizer)

	fileServicesClient := storage.NewFileServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&fileServicesClient.Client, options.ResourceManagerAuthorizer)

	queueServiceClient := storage.NewQueueServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&queueServiceClient.Client, options.ResourceManagerAuthorizer)

	tableServiceClient := storage.NewTableServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&tableServiceClient.Client, options.ResourceManagerAuthorizer)

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	client := Client{
		AccountsClient:              &accountsClient,
		FileSystemsClient:           &fileSystemsClient,
		ADLSGen2PathsClient:         &adlsGen2PathsClient,
		ManagementPoliciesClient:    &managementPoliciesClient,
		BlobInventoryPoliciesClient: &blobInventoryPoliciesClient,
		CloudEndpointsClient:        &cloudEndpointsClient,
		EncryptionScopesClient:      &encryptionScopesClient,
		Environment:                 options.Environment,
		ObjectReplicationClient:     &objectReplicationPolicyClient,
		SubscriptionId:              options.SubscriptionId,
		SyncServiceClient:           &syncServiceClient,
		SyncGroupsClient:            &syncGroupsClient,

		BlobServicesClient:  &blobServicesClient,
		FileServicesClient:  &fileServicesClient,
		QueueServicesClient: &queueServiceClient,
		TableServicesClient: &tableServiceClient,

		resourceManagerAuthorizer: options.ResourceManagerAuthorizer,

		useResourceManager: options.Features.Storage.UseResourceManager,
	}

	if options.StorageUseAzureAD {
		client.storageAdAuth = &options.StorageAuthorizer
	}

	return &client
}

// A setter to manipulate the "useResourceManager" feature flag. This is only meant for the acc test.
func (client *Client) UseResourceManager(t bool) {
	client.useResourceManager = t
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
	if client.useResourceManager {
		rmClient := storage.NewBlobContainersClientWithBaseURI(client.Environment.ResourceManagerEndpoint, client.SubscriptionId)
		rmClient.Client.Authorizer = client.resourceManagerAuthorizer
		return shim.NewManagementPlaneStorageContainerWrapper(&rmClient), nil
	}

	containersClient, err := client.ContainersDataPlaneClient(ctx, account)
	if err != nil {
		return nil, err
	}
	return shim.NewDataPlaneStorageContainerWrapper(containersClient), nil
}

func (client Client) ContainersDataPlaneClient(ctx context.Context, account accountDetails) (*containers.Client, error) {
	if client.storageAdAuth != nil {
		containersClient := containers.NewWithEnvironment(client.Environment)
		containersClient.Client.Authorizer = *client.storageAdAuth
		return &containersClient, nil
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
	return &containersClient, nil
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
	if client.useResourceManager {
		sharesClient := storage.NewFileSharesClientWithBaseURI(client.Environment.ResourceManagerEndpoint, client.SubscriptionId)
		sharesClient.Client.Authorizer = client.resourceManagerAuthorizer
		return shim.NewManagementPlaneStorageShareWrapper(&sharesClient), nil
	}

	sharesClient, err := client.FileSharesDataPlaneClient(ctx, account)
	if err != nil {
		return nil, err
	}
	return shim.NewDataPlaneStorageShareWrapper(sharesClient), nil
}

func (client Client) FileSharesDataPlaneClient(ctx context.Context, account accountDetails) (*shares.Client, error) {
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
	return &sharesClient, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (shim.StorageQueuesWrapper, error) {
	if client.useResourceManager {
		queueClient := storage.NewQueueClient(client.SubscriptionId)
		queueClient.Client.Authorizer = client.resourceManagerAuthorizer
		return shim.NewManagementPlaneStorageQueueWrapper(&queueClient), nil
	}

	queuesClient, err := client.QueuesDataPlaneClient(ctx, account)
	if err != nil {
		return nil, err
	}
	return shim.NewDataPlaneStorageQueueWrapper(queuesClient), nil
}

func (client Client) QueuesDataPlaneClient(ctx context.Context, account accountDetails) (*queues.Client, error) {
	if client.storageAdAuth != nil {
		queueClient := queues.NewWithEnvironment(client.Environment)
		queueClient.Client.Authorizer = *client.storageAdAuth
		return &queueClient, nil
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
	return &queuesClient, nil
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
	// TODO: once mgmt API got ACL support, we can uncomment below
	// if client.useResourceManager {
	//	tableClient := storage.NewTableClient(client.SubscriptionId)
	//	tableClient.Client.Authorizer = client.resourceManagerAuthorizer
	//	return shim.NewManagementPlaneStorageTableWrapper(&tableClient), nil
	// }

	tablesClient, err := client.TablesDataPlaneClient(ctx, account)
	if err != nil {
		return nil, err
	}
	return shim.NewDataPlaneStorageTableWrapper(tablesClient), nil
}

func (client Client) TablesDataPlaneClient(ctx context.Context, account accountDetails) (*tables.Client, error) {
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
	return &tablesClient, nil
}
