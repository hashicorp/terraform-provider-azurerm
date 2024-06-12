// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/cosmosdb"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		CosmosDbMongoUserDefinitionResource{},
		CosmosDbPostgreSQLClusterResource{},
		CosmosDbPostgreSQLCoordinatorConfigurationResource{},
		CosmosDbPostgreSQLFirewallRuleResource{},
		CosmosDbPostgreSQLNodeConfigurationResource{},
		CosmosDbPostgreSQLRoleResource{},
		CosmosDbSqlDedicatedGatewayResource{},
		CosmosDbMongoRoleDefinitionResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "CosmosDB"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"CosmosDB (DocumentDB)",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_cosmosdb_account":                      dataSourceCosmosDbAccount(),
		"azurerm_cosmosdb_mongo_database":               dataSourceCosmosDbMongoDatabase(),
		"azurerm_cosmosdb_restorable_database_accounts": dataSourceCosmosDbRestorableDatabaseAccounts(),
		"azurerm_cosmosdb_sql_database":                 dataSourceCosmosDbSQLDatabase(),
		"azurerm_cosmosdb_sql_role_definition":          dataSourceCosmosDbSQLRoleDefinition(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_cosmosdb_account":              resourceCosmosDbAccount(),
		"azurerm_cosmosdb_cassandra_cluster":    resourceCassandraCluster(),
		"azurerm_cosmosdb_cassandra_datacenter": resourceCassandraDatacenter(),
		"azurerm_cosmosdb_cassandra_keyspace":   resourceCosmosDbCassandraKeyspace(),
		"azurerm_cosmosdb_cassandra_table":      resourceCosmosDbCassandraTable(),
		"azurerm_cosmosdb_gremlin_database":     resourceCosmosGremlinDatabase(),
		"azurerm_cosmosdb_gremlin_graph":        resourceCosmosDbGremlinGraph(),
		"azurerm_cosmosdb_mongo_collection":     resourceCosmosDbMongoCollection(),
		"azurerm_cosmosdb_mongo_database":       resourceCosmosDbMongoDatabase(),
		"azurerm_cosmosdb_sql_container":        resourceCosmosDbSQLContainer(),
		"azurerm_cosmosdb_sql_database":         resourceCosmosDbSQLDatabase(),
		"azurerm_cosmosdb_sql_function":         resourceCosmosDbSQLFunction(),
		"azurerm_cosmosdb_sql_role_assignment":  resourceCosmosDbSQLRoleAssignment(),
		"azurerm_cosmosdb_sql_role_definition":  resourceCosmosDbSQLRoleDefinition(),
		"azurerm_cosmosdb_sql_stored_procedure": resourceCosmosDbSQLStoredProcedure(),
		"azurerm_cosmosdb_sql_trigger":          resourceCosmosDbSQLTrigger(),
		"azurerm_cosmosdb_table":                resourceCosmosDbTable(),
	}

	if !features.FourPointOhBeta() {
		resources["azurerm_cosmosdb_notebook_workspace"] = resourceCosmosDbNotebookWorkspace()
	}

	return resources
}
