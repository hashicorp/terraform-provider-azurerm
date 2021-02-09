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
		"azurerm_kubernetes_service_versions":  dataSourceKubernetesServiceVersions(),
		"azurerm_container_registry":           dataSourceContainerRegistry(),
		"azurerm_kubernetes_cluster":           dataSourceKubernetesCluster(),
		"azurerm_kubernetes_cluster_node_pool": dataSourceKubernetesClusterNodePool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_container_group":              resourceContainerGroup(),
		"azurerm_container_registry_webhook":   resourceContainerRegistryWebhook(),
		"azurerm_container_registry":           resourceContainerRegistry(),
		"azurerm_kubernetes_cluster":           resourceKubernetesCluster(),
		"azurerm_kubernetes_cluster_node_pool": resourceKubernetesClusterNodePool(),
	}
}
