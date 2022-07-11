package databricks

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
	return map[string]*pluginsdk.Resource{
		"azurerm_databricks_workspace":                      resourceDatabricksWorkspace(),
		"azurerm_databricks_workspace_customer_managed_key": resourceDatabricksWorkspaceCustomerManagedKey(),
	}
}
