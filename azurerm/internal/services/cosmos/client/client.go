package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CassandraClient *documentdb.CassandraResourcesClient
	DatabaseClient  *documentdb.DatabaseAccountsClient
	GremlinClient   *documentdb.GremlinResourcesClient
	MongoDbClient   *documentdb.MongoDBResourcesClient
	SqlClient       *documentdb.SQLResourcesClient
	TableClient     *documentdb.TableResourcesClient
}

func NewClient(o *common.ClientOptions) *Client {
	cassandraClient := documentdb.NewCassandraResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cassandraClient.Client, o.ResourceManagerAuthorizer)

	databaseClient := documentdb.NewDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	gremlinClient := documentdb.NewGremlinResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gremlinClient.Client, o.ResourceManagerAuthorizer)

	mongoDbClient := documentdb.NewMongoDBResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&mongoDbClient.Client, o.ResourceManagerAuthorizer)

	sqlClient := documentdb.NewSQLResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlClient.Client, o.ResourceManagerAuthorizer)

	tableClient := documentdb.NewTableResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tableClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CassandraClient: &cassandraClient,
		DatabaseClient:  &databaseClient,
		GremlinClient:   &gremlinClient,
		MongoDbClient:   &mongoDbClient,
		SqlClient:       &sqlClient,
		TableClient:     &tableClient,
	}
}
