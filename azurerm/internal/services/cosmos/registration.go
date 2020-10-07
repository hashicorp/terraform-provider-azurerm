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
		"azurerm_cosmosdb_account": dataSourceArmCosmosDbAccount(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_cosmosdb_account":              resourceArmCosmosDbAccount(),
		"azurerm_cosmosdb_cassandra_keyspace":   resourceArmCosmosDbCassandraKeyspace(),
		"azurerm_cosmosdb_gremlin_database":     resourceArmCosmosGremlinDatabase(),
		"azurerm_cosmosdb_gremlin_graph":        resourceArmCosmosDbGremlinGraph(),
		"azurerm_cosmosdb_mongo_collection":     resourceArmCosmosDbMongoCollection(),
		"azurerm_cosmosdb_mongo_database":       resourceArmCosmosDbMongoDatabase(),
		"azurerm_cosmosdb_sql_container":        resourceArmCosmosDbSQLContainer(),
		"azurerm_cosmosdb_sql_database":         resourceArmCosmosDbSQLDatabase(),
		"azurerm_cosmosdb_sql_stored_procedure": resourceArmCosmosDbSQLStoredProcedure(),
		"azurerm_cosmosdb_table":                resourceArmCosmosDbTable(),
	}
}
