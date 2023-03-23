package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2022-01-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Client         *redisenterprise.RedisEnterpriseClient
	DatabaseClient *databases.DatabasesClient
}

func NewClient(o *common.ClientOptions) *Client {
	client := redisenterprise.NewRedisEnterpriseClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	databaseClient := databases.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client:         &client,
		DatabaseClient: &databaseClient,
	}
}
