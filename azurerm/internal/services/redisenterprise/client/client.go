package client

import (
	"github.com/Azure/azure-sdk-for-go/services/redisenterprise/mgmt/2021-03-01/redisenterprise"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client           *redisenterprise.Client
	DatabaseClient   *redisenterprise.DatabasesClient
	OperationsClient *redisenterprise.OperationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	client := redisenterprise.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	databaseClient := redisenterprise.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	operationsClient := redisenterprise.NewOperationsClient(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client:           &client,
		DatabaseClient:   &databaseClient,
		OperationsClient: &operationsClient,
	}
}
