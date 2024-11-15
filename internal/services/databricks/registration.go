// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

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
		"azurerm_databricks_workspace":                             dataSourceDatabricksWorkspace(),
		"azurerm_databricks_workspace_private_endpoint_connection": dataSourceDatabricksWorkspacePrivateEndpointConnection(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_databricks_workspace":                                resourceDatabricksWorkspace(),
		"azurerm_databricks_workspace_root_dbfs_customer_managed_key": resourceDatabricksWorkspaceRootDbfsCustomerManagedKey(),
		"azurerm_databricks_virtual_network_peering":                  resourceDatabricksVirtualNetworkPeering(),
	}

	if !features.FivePointOhBeta() {
		resources["azurerm_databricks_workspace_customer_managed_key"] = resourceDatabricksWorkspaceCustomerManagedKey()
	}

	return resources
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		DatabricksAccessConnectorDataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AccessConnectorResource{},
	}
}
