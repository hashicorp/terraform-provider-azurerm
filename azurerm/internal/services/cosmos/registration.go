package cosmos

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_cosmosdb_account": dataSourceCosmosDbAccount(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_cosmosdb_account":              resourceCosmosDbAccount(),
		"azurerm_cosmosdb_cassandra_keyspace":   resourceCosmosDbCassandraKeyspace(),
		"azurerm_cosmosdb_cassandra_table":      resourceCosmosDbCassandraTable(),
		"azurerm_cosmosdb_gremlin_database":     resourceCosmosGremlinDatabase(),
		"azurerm_cosmosdb_gremlin_graph":        resourceCosmosDbGremlinGraph(),
		"azurerm_cosmosdb_mongo_collection":     resourceCosmosDbMongoCollection(),
		"azurerm_cosmosdb_mongo_database":       resourceCosmosDbMongoDatabase(),
		"azurerm_cosmosdb_sql_container":        resourceCosmosDbSQLContainer(),
		"azurerm_cosmosdb_sql_database":         resourceCosmosDbSQLDatabase(),
		"azurerm_cosmosdb_sql_stored_procedure": resourceCosmosDbSQLStoredProcedure(),
		"azurerm_cosmosdb_table":                resourceCosmosDbTable(),
	}
}
