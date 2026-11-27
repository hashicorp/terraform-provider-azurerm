// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_cassandra_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_cassandra_datacenter"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_cassandra_keyspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_cassandra_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_gremlin_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_gremlin_graph"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_mongo_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_mongo_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_mongo_role_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_mongo_user_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_postgresql_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_postgresql_coordinator_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_postgresql_firewall_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_postgresql_node_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_postgresql_role"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_restorable_database_accounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_container"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_dedicated_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_function"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_role_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_stored_procedure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_sql_trigger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/resources/cosmosdb_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
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
		cosmosdb_mongo_role_definition.Resource{},
		cosmosdb_mongo_user_definition.Resource{},
		cosmosdb_postgresql_cluster.Resource{},
		cosmosdb_postgresql_coordinator_configuration.Resource{},
		cosmosdb_postgresql_firewall_rule.Resource{},
		cosmosdb_postgresql_node_configuration.Resource{},
		cosmosdb_postgresql_role.Resource{},
		cosmosdb_sql_dedicated_gateway.Resource{},
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
		"azurerm_cosmosdb_account":                      cosmosdb_account.RegisterDataSource(),
		"azurerm_cosmosdb_mongo_database":               cosmosdb_mongo_database.RegisterDataSource(),
		"azurerm_cosmosdb_restorable_database_accounts": cosmosdb_restorable_database_accounts.RegisterDataSource(),
		"azurerm_cosmosdb_sql_database":                 cosmosdb_sql_database.RegisterDataSource(),
		"azurerm_cosmosdb_sql_role_definition":          cosmosdb_sql_role_definition.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_cosmosdb_account":              cosmosdb_account.RegisterResource(),
		"azurerm_cosmosdb_cassandra_cluster":    cosmosdb_cassandra_cluster.RegisterResource(),
		"azurerm_cosmosdb_cassandra_datacenter": cosmosdb_cassandra_datacenter.RegisterResource(),
		"azurerm_cosmosdb_cassandra_keyspace":   cosmosdb_cassandra_keyspace.RegisterResource(),
		"azurerm_cosmosdb_cassandra_table":      cosmosdb_cassandra_table.RegisterResource(),
		"azurerm_cosmosdb_gremlin_database":     cosmosdb_gremlin_database.RegisterResource(),
		"azurerm_cosmosdb_gremlin_graph":        cosmosdb_gremlin_graph.RegisterResource(),
		"azurerm_cosmosdb_mongo_collection":     cosmosdb_mongo_collection.RegisterResource(),
		"azurerm_cosmosdb_mongo_database":       cosmosdb_mongo_database.RegisterResource(),
		"azurerm_cosmosdb_sql_container":        cosmosdb_sql_container.RegisterResource(),
		"azurerm_cosmosdb_sql_database":         cosmosdb_sql_database.RegisterResource(),
		"azurerm_cosmosdb_sql_function":         cosmosdb_sql_function.RegisterResource(),
		"azurerm_cosmosdb_sql_role_assignment":  cosmosdb_sql_role_assignment.RegisterResource(),
		"azurerm_cosmosdb_sql_role_definition":  cosmosdb_sql_role_definition.RegisterResource(),
		"azurerm_cosmosdb_sql_stored_procedure": cosmosdb_sql_stored_procedure.RegisterResource(),
		"azurerm_cosmosdb_sql_trigger":          cosmosdb_sql_trigger.RegisterResource(),
		"azurerm_cosmosdb_table":                cosmosdb_table.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
