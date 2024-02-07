// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-05-15/sqldedicatedgateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2022-11-15/mongorbacs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2023-04-15/managedcassandras"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/roles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CassandraClient                  *documentdb.CassandraResourcesClient
	ClustersClient                   *clusters.ClustersClient
	ConfigurationsClient             *configurations.ConfigurationsClient
	CosmosDBClient                   *cosmosdb.CosmosDBClient
	DatabaseClient                   *documentdb.DatabaseAccountsClient
	FirewallRulesClient              *firewallrules.FirewallRulesClient
	GremlinClient                    *documentdb.GremlinResourcesClient
	ManagedCassandraClient           *managedcassandras.ManagedCassandrasClient
	MongoDbClient                    *documentdb.MongoDBResourcesClient
	MongoRBACClient                  *mongorbacs.MongorbacsClient
	NotebookWorkspaceClient          *documentdb.NotebookWorkspacesClient
	RestorableDatabaseAccountsClient *documentdb.RestorableDatabaseAccountsClient
	RolesClient                      *roles.RolesClient
	SqlDedicatedGatewayClient        *sqldedicatedgateway.SqlDedicatedGatewayClient
	SqlClient                        *documentdb.SQLResourcesClient
	SqlResourceClient                *documentdb.SQLResourcesClient
	TableClient                      *documentdb.TableResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	cassandraClient := documentdb.NewCassandraResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&cassandraClient.Client, o.ResourceManagerAuthorizer)

	managedCassandraClient := managedcassandras.NewManagedCassandrasClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&managedCassandraClient.Client, o.ResourceManagerAuthorizer)

	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	configurationsClient, err := configurations.NewConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Configurations client: %+v", err)
	}
	o.Configure(configurationsClient.Client, o.Authorizers.ResourceManager)

	cosmosdbClient := cosmosdb.NewCosmosDBClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&cosmosdbClient.Client, o.ResourceManagerAuthorizer)

	databaseClient := documentdb.NewDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FirewallRules client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	gremlinClient := documentdb.NewGremlinResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&gremlinClient.Client, o.ResourceManagerAuthorizer)

	mongoDbClient := documentdb.NewMongoDBResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&mongoDbClient.Client, o.ResourceManagerAuthorizer)

	mongorbacsClient := mongorbacs.NewMongorbacsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&mongorbacsClient.Client, o.ResourceManagerAuthorizer)

	notebookWorkspaceClient := documentdb.NewNotebookWorkspacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&notebookWorkspaceClient.Client, o.ResourceManagerAuthorizer)

	restorableDatabaseAccountsClient := documentdb.NewRestorableDatabaseAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&restorableDatabaseAccountsClient.Client, o.ResourceManagerAuthorizer)

	rolesClient, err := roles.NewRolesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Roles client: %+v", err)
	}
	o.Configure(rolesClient.Client, o.Authorizers.ResourceManager)

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
		ManagedCassandraClient:           &managedCassandraClient,
		ClustersClient:                   clustersClient,
		ConfigurationsClient:             configurationsClient,
		CosmosDBClient:                   &cosmosdbClient,
		DatabaseClient:                   &databaseClient,
		FirewallRulesClient:              firewallRulesClient,
		GremlinClient:                    &gremlinClient,
		MongoDbClient:                    &mongoDbClient,
		MongoRBACClient:                  &mongorbacsClient,
		NotebookWorkspaceClient:          &notebookWorkspaceClient,
		RestorableDatabaseAccountsClient: &restorableDatabaseAccountsClient,
		RolesClient:                      rolesClient,
		SqlDedicatedGatewayClient:        &sqlDedicatedGatewayClient,
		SqlClient:                        &sqlClient,
		SqlResourceClient:                &sqlResourceClient,
		TableClient:                      &tableClient,
	}, nil
}
