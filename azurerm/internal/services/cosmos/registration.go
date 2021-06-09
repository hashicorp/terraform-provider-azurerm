package cosmos

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_cosmosdb_account": dataSourceCosmosDbAccount(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_cosmosdb_account":              resourceCosmosDbAccount(),
		"azurerm_cosmosdb_cassandra_keyspace":   resourceCosmosDbCassandraKeyspace(),
		"azurerm_cosmosdb_cassandra_table":      resourceCosmosDbCassandraTable(),
		"azurerm_cosmosdb_gremlin_database":     resourceCosmosGremlinDatabase(),
		"azurerm_cosmosdb_gremlin_graph":        resourceCosmosDbGremlinGraph(),
		"azurerm_cosmosdb_mongo_collection":     resourceCosmosDbMongoCollection(),
		"azurerm_cosmosdb_mongo_database":       resourceCosmosDbMongoDatabase(),
		"azurerm_cosmosdb_notebook_workspace":   resourceCosmosDbNotebookWorkspace(),
		"azurerm_cosmosdb_sql_container":        resourceCosmosDbSQLContainer(),
		"azurerm_cosmosdb_sql_database":         resourceCosmosDbSQLDatabase(),
		"azurerm_cosmosdb_sql_function":         resourceCosmosDbSQLFunction(),
		"azurerm_cosmosdb_sql_stored_procedure": resourceCosmosDbSQLStoredProcedure(),
		"azurerm_cosmosdb_sql_trigger":          resourceCosmosDbSQLTrigger(),
		"azurerm_cosmosdb_table":                resourceCosmosDbTable(),
	}
}
