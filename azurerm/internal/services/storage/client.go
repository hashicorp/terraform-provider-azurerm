package storage

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/authorizers"
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
	AccountsClient storage.AccountsClient

	environment az.Environment
}

func BuildClient(options *common.ClientOptions) *Client {
	accountsClient := storage.NewAccountsClientWithBaseURI(options.ResourceManagerEndpoint, options.SubscriptionId)
	options.ConfigureClient(&accountsClient.Client, options.ResourceManagerAuthorizer)

	// TODO: switch Storage Containers to using the storage.BlobContainersClient
	// (which should fix #2977) when the storage clients have been moved in here
	return &Client{
		AccountsClient: accountsClient,
		environment:    options.Environment,
	}
}

func (client Client) BlobsClient(ctx context.Context, resourceGroup, accountName string) (*blobs.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteAuthorizer(accountName, *accountKey)
	blobsClient := blobs.NewWithEnvironment(client.environment)
	blobsClient.Client.Authorizer = storageAuth
	return &blobsClient, nil
}

func (client Client) ContainersClient(ctx context.Context, resourceGroup, accountName string) (*containers.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteAuthorizer(accountName, *accountKey)
	containersClient := containers.NewWithEnvironment(client.environment)
	containersClient.Client.Authorizer = storageAuth
	return &containersClient, nil
}

func (client Client) FileShareDirectoriesClient(ctx context.Context, resourceGroup, accountName string) (*directories.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteAuthorizer(accountName, *accountKey)
	directoriesClient := directories.NewWithEnvironment(client.environment)
	directoriesClient.Client.Authorizer = storageAuth
	return &directoriesClient, nil
}

func (client Client) FileSharesClient(ctx context.Context, resourceGroup, accountName string) (*shares.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteAuthorizer(accountName, *accountKey)
	directoriesClient := shares.NewWithEnvironment(client.environment)
	directoriesClient.Client.Authorizer = storageAuth
	return &directoriesClient, nil
}

func (client Client) QueuesClient(ctx context.Context, resourceGroup, accountName string) (*queues.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteAuthorizer(accountName, *accountKey)
	queuesClient := queues.NewWithEnvironment(client.environment)
	queuesClient.Client.Authorizer = storageAuth
	return &queuesClient, nil
}

func (client Client) TableEntityClient(ctx context.Context, resourceGroup, accountName string) (*entities.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteTableAuthorizer(accountName, *accountKey)
	entitiesClient := entities.NewWithEnvironment(client.environment)
	entitiesClient.Client.Authorizer = storageAuth
	return &entitiesClient, nil
}

func (client Client) TablesClient(ctx context.Context, resourceGroup, accountName string) (*tables.Client, error) {
	accountKey, err := client.findAccountKey(ctx, resourceGroup, accountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Account Key: %s", err)
	}

	storageAuth := authorizers.NewSharedKeyLiteTableAuthorizer(accountName, *accountKey)
	tablesClient := tables.NewWithEnvironment(client.environment)
	tablesClient.Client.Authorizer = storageAuth
	return &tablesClient, nil
}
