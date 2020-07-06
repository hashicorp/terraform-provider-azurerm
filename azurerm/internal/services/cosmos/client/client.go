package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2020-04-01/documentdb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	CassandraResourcesClient *documentdb.CassandraResourcesClient
	DatabaseClient           *documentdb.DatabaseAccountsClient
	GremlinResourcesClient   *documentdb.GremlinResourcesClient
	MongoDBResourcesClient   *documentdb.MongoDBResourcesClient
	SQLResourcesClient       *documentdb.SQLResourcesClient
	TableResourcesClient     *documentdb.TableResourcesClient
}

func NewClient(o *common.ClientOptions) *Client {
	cassandraResourcesClient := documentdb.NewCassandraResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cassandraResourcesClient.Client, o.ResourceManagerAuthorizer)

	databaseClient := documentdb.NewDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	gremlinResourcesClient := documentdb.NewGremlinResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gremlinResourcesClient.Client, o.ResourceManagerAuthorizer)

	mongoDBResourcesClient := documentdb.NewMongoDBResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&mongoDBResourcesClient.Client, o.ResourceManagerAuthorizer)

	sqlResourcesClient := documentdb.NewSQLResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlResourcesClient.Client, o.ResourceManagerAuthorizer)

	tableResourcesClient := documentdb.NewTableResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tableResourcesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CassandraResourcesClient: &cassandraResourcesClient,
		DatabaseClient:           &databaseClient,
		GremlinResourcesClient:   &gremlinResourcesClient,
		MongoDBResourcesClient:   &mongoDBResourcesClient,
		SQLResourcesClient:       &sqlResourcesClient,
		TableResourcesClient:     &tableResourcesClient,
	}
}
