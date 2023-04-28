package client

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/managedcassandras"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CassandraClient                  *documentdb.CassandraResourcesClient
	CassandraClustersClient          *managedcassandras.ManagedCassandrasClient
	CassandraDatacentersClient       *documentdb.CassandraDataCentersClient
	ClustersClient                   *clusters.ClustersClient
	ConfigurationsClient             *configurations.ConfigurationsClient
	DatabaseClient                   *documentdb.DatabaseAccountsClient
	GremlinClient                    *documentdb.GremlinResourcesClient
	MongoDbClient                    *documentdb.MongoDBResourcesClient
	NotebookWorkspaceClient          *documentdb.NotebookWorkspacesClient
	RestorableDatabaseAccountsClient *documentdb.RestorableDatabaseAccountsClient
	SqlDedicatedGatewayClient        *sqldedicatedgateway.SqlDedicatedGatewayClient
	SqlClient                        *documentdb.SQLResourcesClient
	SqlResourceClient                *documentdb.SQLResourcesClient
	TableClient                      *documentdb.TableResourcesClient
}

func NewClient(o *common.ClientOptions) *Client {
	cassandraClient := documentdb.NewCassandraResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cassandraClient.Client, o.ResourceManagerAuthorizer)

	cassandraClustersClient := managedcassandras.NewManagedCassandrasClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&cassandraClustersClient.Client, o.ResourceManagerAuthorizer)

	cassandraDatacentersClient := documentdb.NewCassandraDataCentersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cassandraDatacentersClient.Client, o.ResourceManagerAuthorizer)

	clustersClient := clusters.NewClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clustersClient.Client, o.ResourceManagerAuthorizer)

	configurationsClient := configurations.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	databaseClient := documentdb.NewDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	gremlinClient := documentdb.NewGremlinResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gremlinClient.Client, o.ResourceManagerAuthorizer)

	mongoDbClient := documentdb.NewMongoDBResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&mongoDbClient.Client, o.ResourceManagerAuthorizer)

	notebookWorkspaceClient := documentdb.NewNotebookWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&notebookWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	restorableDatabaseAccountsClient := documentdb.NewRestorableDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&restorableDatabaseAccountsClient.Client, o.ResourceManagerAuthorizer)

	sqlDedicatedGatewayClient := sqldedicatedgateway.NewSqlDedicatedGatewayClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&sqlDedicatedGatewayClient.Client, o.ResourceManagerAuthorizer)

	sqlClient := documentdb.NewSQLResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlClient.Client, o.ResourceManagerAuthorizer)

	sqlResourceClient := documentdb.NewSQLResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlResourceClient.Client, o.ResourceManagerAuthorizer)

	tableClient := documentdb.NewTableResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tableClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		CassandraClient:                  &cassandraClient,
		CassandraClustersClient:          &cassandraClustersClient,
		CassandraDatacentersClient:       &cassandraDatacentersClient,
		ClustersClient:                   &clustersClient,
		ConfigurationsClient:             &configurationsClient,
		DatabaseClient:                   &databaseClient,
		GremlinClient:                    &gremlinClient,
		MongoDbClient:                    &mongoDbClient,
		NotebookWorkspaceClient:          &notebookWorkspaceClient,
		RestorableDatabaseAccountsClient: &restorableDatabaseAccountsClient,
		SqlDedicatedGatewayClient:        &sqlDedicatedGatewayClient,
		SqlClient:                        &sqlClient,
		SqlResourceClient:                &sqlResourceClient,
		TableClient:                      &tableClient,
	}
}
