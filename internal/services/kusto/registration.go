// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_attached_database_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_cluster_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_cluster_managed_private_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_cluster_principal_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_cosmosdb_data_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_database_principal_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_eventgrid_data_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_eventhub_data_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_iothub_data_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/resources/kusto_script"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		kusto_cosmosdb_data_connection.Resource{},
	}
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/kusto"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Kusto"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Data Explorer",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_kusto_cluster":  kusto_cluster.RegisterDataSource(),
		"azurerm_kusto_database": kusto_database.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_kusto_attached_database_configuration":  kusto_attached_database_configuration.RegisterResource(),
		"azurerm_kusto_cluster":                          kusto_cluster.RegisterResource(),
		"azurerm_kusto_cluster_customer_managed_key":     kusto_cluster_customer_managed_key.RegisterResource(),
		"azurerm_kusto_cluster_managed_private_endpoint": kusto_cluster_managed_private_endpoint.RegisterResource(),
		"azurerm_kusto_cluster_principal_assignment":     kusto_cluster_principal_assignment.RegisterResource(),
		"azurerm_kusto_database":                         kusto_database.RegisterResource(),
		"azurerm_kusto_database_principal_assignment":    kusto_database_principal_assignment.RegisterResource(),
		"azurerm_kusto_eventgrid_data_connection":        kusto_eventgrid_data_connection.RegisterResource(),
		"azurerm_kusto_eventhub_data_connection":         kusto_eventhub_data_connection.RegisterResource(),
		"azurerm_kusto_iothub_data_connection":           kusto_iothub_data_connection.RegisterResource(),
		"azurerm_kusto_script":                           kusto_script.RegisterResource(),
	}
}
