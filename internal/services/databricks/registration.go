// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/resources/databricks_access_connector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/resources/databricks_virtual_network_peering"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/resources/databricks_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/resources/databricks_workspace_private_endpoint_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/resources/databricks_workspace_root_dbfs_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/databricks"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "DataBricks"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Databricks",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_databricks_workspace":                             databricks_workspace.RegisterDataSource(),
		"azurerm_databricks_workspace_private_endpoint_connection": databricks_workspace_private_endpoint_connection.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_databricks_virtual_network_peering":                  databricks_virtual_network_peering.RegisterResource(),
		"azurerm_databricks_workspace":                                databricks_workspace.RegisterResource(),
		"azurerm_databricks_workspace_root_dbfs_customer_managed_key": databricks_workspace_root_dbfs_customer_managed_key.RegisterResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		databricks_access_connector.DataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		databricks_access_connector.Resource{},
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
