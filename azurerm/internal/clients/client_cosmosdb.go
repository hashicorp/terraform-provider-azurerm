package clients

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type CosmosDBClient struct {
	DatabaseClient *documentdb.DatabaseAccountsClient
}

func newCosmosDBClient(o *common.ClientOptions) *CosmosDBClient {
	databaseClient := documentdb.NewDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	return &CosmosDBClient{
		DatabaseClient: &databaseClient,
	}
}
