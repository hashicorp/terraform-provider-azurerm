package containers

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Container Services"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Container",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_kubernetes_service_versions":  dataSourceArmKubernetesServiceVersions(),
		"azurerm_container_registry":           dataSourceArmContainerRegistry(),
		"azurerm_container_registry_token":     dataSourceArmContainerRegistryToken(),
		"azurerm_container_registry_scope_map": dataSourceArmContainerRegistryScopeMap(),
		"azurerm_kubernetes_cluster":           dataSourceArmKubernetesCluster(),
		"azurerm_kubernetes_cluster_node_pool": dataSourceKubernetesClusterNodePool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_container_group":              resourceArmContainerGroup(),
		"azurerm_container_registry_webhook":   resourceArmContainerRegistryWebhook(),
		"azurerm_container_registry":           resourceArmContainerRegistry(),
		"azurerm_container_registry_token":     resourceArmContainerRegistryToken(),
		"azurerm_container_registry_scope_map": resourceArmContainerRegistryScopeMap(),
		"azurerm_kubernetes_cluster":           resourceArmKubernetesCluster(),
		"azurerm_kubernetes_cluster_node_pool": resourceArmKubernetesClusterNodePool(),
	}
}
