package synapse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Synapse"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Synapse",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_synapse_workspace": dataSourceSynapseWorkspace(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_synapse_firewall_rule":            resourceSynapseFirewallRule(),
		"azurerm_synapse_private_link_hub":         resourceSynapsePrivateLinkHub(),
		"azurerm_synapse_managed_private_endpoint": resourceSynapseManagedPrivateEndpoint(),
		"azurerm_synapse_role_assignment":          resourceSynapseRoleAssignment(),
		"azurerm_synapse_spark_pool":               resourceSynapseSparkPool(),
		"azurerm_synapse_sql_pool":                 resourceSynapseSqlPool(),
		"azurerm_synapse_workspace":                resourceSynapseWorkspace(),
	}
}
