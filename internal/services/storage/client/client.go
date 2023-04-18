package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"         // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2020-03-01/storagesync" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	storage_v2022_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2022-05-01/localusers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
	LocalUsersClient            *localusers.LocalUsersClient
	FileSystemsClient           *filesystems.Client
	ADLSGen2PathsClient         *paths.Client
	ManagementPoliciesClient    *storage.ManagementPoliciesClient
	BlobServicesClient          *storage.BlobServicesClient
	BlobInventoryPoliciesClient *storage.BlobInventoryPoliciesClient
	CloudEndpointsClient        *storagesync.CloudEndpointsClient
	EncryptionScopesClient      *storage.EncryptionScopesClient
	Environment                 azure.Environment
	FileServicesClient          *storage.FileServicesClient
	SyncServiceClient           *storagesync.ServicesClient
	SyncGroupsClient            *storagesync.SyncGroupsClient
	SubscriptionId              string

	// Data plane clients using AAD auth
	AccountsDataPlaneAADClient   *accounts.Client
	BlobsDataPlaneAADClient      *blobs.Client
	ContainersDataPlaneAADClient *containers.Client
	QueuesDataPlaneAADClient     *queues.Client

	ResourceManager *storage_v2022_05_01.Client

	resourceManagerAuthorizer autorest.Authorizer
	storageAdAuth             *autorest.Authorizer

	// The client option is used to configure on-demand data plane clients
	clientOpt *common.ClientOptions
}

func NewClient(options *common.ClientOptions) *Client {
	accountsClient := storage.NewAccountsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&accountsClient.Client, options.ResourceManagerAuthorizer)

	localUsersClient := localusers.NewLocalUsersClientWithBaseURI(options.ResourceManagerEndpoint)
	localUsersClient.Client.Authorizer = options.ResourceManagerAuthorizer

	fileSystemsClient := filesystems.NewWithEnvironment(options.AzureEnvironment)
	options.ConfigureClient(&fileSystemsClient.Client, options.StorageAuthorizer)

	adlsGen2PathsClient := paths.NewWithEnvironment(options.AzureEnvironment)
	options.ConfigureClient(&adlsGen2PathsClient.Client, options.StorageAuthorizer)

	managementPoliciesClient := storage.NewManagementPoliciesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&managementPoliciesClient.Client, options.ResourceManagerAuthorizer)

	blobServicesClient := storage.NewBlobServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&blobServicesClient.Client, options.ResourceManagerAuthorizer)

	blobInventoryPoliciesClient := storage.NewBlobInventoryPoliciesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&blobInventoryPoliciesClient.Client, options.ResourceManagerAuthorizer)

	cloudEndpointsClient := storagesync.NewCloudEndpointsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&cloudEndpointsClient.Client, options.ResourceManagerAuthorizer)

	encryptionScopesClient := storage.NewEncryptionScopesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&encryptionScopesClient.Client, options.ResourceManagerAuthorizer)

	fileServicesClient := storage.NewFileServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&fileServicesClient.Client, options.ResourceManagerAuthorizer)

	resourceManager := storage_v2022_05_01.NewClientWithBaseURI(options.ResourceManagerEndpoint,
		func(c *autorest.Client) {
			c.Authorizer = options.ResourceManagerAuthorizer
		})

	syncServiceClient := storagesync.NewServicesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncServiceClient.Client, options.ResourceManagerAuthorizer)

	syncGroupsClient := storagesync.NewSyncGroupsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&syncGroupsClient.Client, options.ResourceManagerAuthorizer)

	accountsDataPlaneAADClient := accounts.NewWithEnvironment(options.AzureEnvironment)
	options.ConfigureClient(&accountsDataPlaneAADClient.Client, options.ResourceManagerAuthorizer)

	blobsDataPlaneAADClient := blobs.NewWithEnvironment(options.AzureEnvironment)
	options.ConfigureClient(&blobsDataPlaneAADClient.Client, options.ResourceManagerAuthorizer)

	containersDataPlaneAADClient := containers.NewWithEnvironment(options.AzureEnvironment)
	options.ConfigureClient(&containersDataPlaneAADClient.Client, options.ResourceManagerAuthorizer)

	queuesDataPlaneAADClient := queues.NewWithEnvironment(options.AzureEnvironment)
	options.ConfigureClient(&queuesDataPlaneAADClient.Client, options.ResourceManagerAuthorizer)

	opts := options.Clone()

	// This is necessary since the data plane client that uses shared key will uses any header values that starts with "x-ms-"
	// to build key signature: https://github.com/hashicorp/terraform-provider-azurerm/blob/8c7f98f1a1efc4c033a5d33e3b2006ef81faf5c6/vendor/github.com/Azure/go-autorest/autorest/authorization_storage.go#L265-L270
	// The correlation request ID is added in the header after the authorization middleware (i.e. the shared key signature calc), see: https://github.com/hashicorp/terraform-provider-azurerm/blob/8c7f98f1a1efc4c033a5d33e3b2006ef81faf5c6/vendor/github.com/Azure/go-autorest/autorest/client.go#L242-L244
	// Therefore, we can't send the request along with this header.
	opts.DisableCorrelationRequestID = true

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	client := Client{
		AccountsClient:              &accountsClient,
		LocalUsersClient:            &localUsersClient,
		FileSystemsClient:           &fileSystemsClient,
		ADLSGen2PathsClient:         &adlsGen2PathsClient,
		ManagementPoliciesClient:    &managementPoliciesClient,
		BlobServicesClient:          &blobServicesClient,
		BlobInventoryPoliciesClient: &blobInventoryPoliciesClient,
		CloudEndpointsClient:        &cloudEndpointsClient,
		EncryptionScopesClient:      &encryptionScopesClient,
		Environment:                 options.AzureEnvironment,
		FileServicesClient:          &fileServicesClient,
		ResourceManager:             &resourceManager,
		SubscriptionId:              options.SubscriptionId,
		SyncServiceClient:           &syncServiceClient,
		SyncGroupsClient:            &syncGroupsClient,

		AccountsDataPlaneAADClient:   &accountsDataPlaneAADClient,
		BlobsDataPlaneAADClient:      &blobsDataPlaneAADClient,
		ContainersDataPlaneAADClient: &containersDataPlaneAADClient,
		QueuesDataPlaneAADClient:     &queuesDataPlaneAADClient,

		resourceManagerAuthorizer: options.ResourceManagerAuthorizer,
		clientOpt:                 &opts,
	}

	if options.StorageUseAzureAD {
		client.storageAdAuth = &options.StorageAuthorizer
	}

	return &client
}

func (client Client) AccountsDataPlaneClient(ctx context.Context, account accountDetails) (*accounts.Client, error) {
	if client.storageAdAuth != nil {
		return client.AccountsDataPlaneAADClient, nil
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
	client.clientOpt.ConfigureClient(&accountsClient.Client, storageAuth)
	return &accountsClient, nil
}

func (client Client) BlobsClient(ctx context.Context, account accountDetails) (*blobs.Client, error) {
	if client.storageAdAuth != nil {
		return client.BlobsDataPlaneAADClient, nil
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
	client.clientOpt.ConfigureClient(&blobsClient.Client, storageAuth)
	return &blobsClient, nil
}

func (client Client) ContainersClient(ctx context.Context, account accountDetails) (shim.StorageContainerWrapper, error) {
	if client.storageAdAuth != nil {
		shim := shim.NewDataPlaneStorageContainerWrapper(client.ContainersDataPlaneAADClient)
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
	client.clientOpt.ConfigureClient(&containersClient.Client, storageAuth)

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
	client.clientOpt.ConfigureClient(&directoriesClient.Client, storageAuth)
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
	client.clientOpt.ConfigureClient(&filesClient.Client, storageAuth)
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
	client.clientOpt.ConfigureClient(&sharesClient.Client, storageAuth)
	shim := shim.NewDataPlaneStorageShareWrapper(&sharesClient)
	return shim, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (shim.StorageQueuesWrapper, error) {
	if client.storageAdAuth != nil {
		return shim.NewDataPlaneStorageQueueWrapper(client.QueuesDataPlaneAADClient), nil
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
	client.clientOpt.ConfigureClient(&queuesClient.Client, storageAuth)
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
	client.clientOpt.ConfigureClient(&entitiesClient.Client, storageAuth)
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
	client.clientOpt.ConfigureClient(&tablesClient.Client, storageAuth)
	shim := shim.NewDataPlaneStorageTableWrapper(&tablesClient)
	return shim, nil
}
