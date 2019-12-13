package client

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/datalakestore/filesystems"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/blobs"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/directories"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/queue/queues"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/entities"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/table/tables"
)

type Client struct {
	AccountsClient           *storage.AccountsClient
	FileSystemsClient        *filesystems.Client
	ManagementPoliciesClient storage.ManagementPoliciesClient

	environment az.Environment
}

func NewClient(options *common.ClientOptions) *Client {
	accountsClient := storage.NewAccountsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&accountsClient.Client, options.ResourceManagerAuthorizer)

	fileSystemsClient := filesystems.NewWithEnvironment(options.Environment)
	fileSystemsClient.Authorizer = options.StorageAuthorizer

	managementPoliciesClient := storage.NewManagementPoliciesClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&managementPoliciesClient.Client, options.ResourceManagerAuthorizer)

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	return &Client{
		AccountsClient:           &accountsClient,
		FileSystemsClient:        &fileSystemsClient,
		ManagementPoliciesClient: managementPoliciesClient,
		environment:              options.Environment,
	}
}

func (client Client) BlobsClient(ctx context.Context, account accountDetails) (*blobs.Client, error) {
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
	accountKey, err := account.AccountKey(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth, err := autorest.NewSharedKeyAuthorizer(account.name, *accountKey, autorest.SharedKeyLite)
	if err != nil {
		return nil, fmt.Errorf("Error building Authorizer: %+v", err)
	}

	directoriesClient := shares.NewWithEnvironment(client.environment)
	directoriesClient.Client.Authorizer = storageAuth
	return &directoriesClient, nil
}

func (client Client) QueuesClient(ctx context.Context, account accountDetails) (*queues.Client, error) {
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
