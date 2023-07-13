// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct {
	autoRegistration
}

var (
	_ sdk.TypedServiceRegistration   = Registration{}
	_ sdk.UntypedServiceRegistration = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "Container Services"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	categories := []string{
		"Container",
	}
	categories = append(categories, r.autoRegistration.WebsiteCategories()...)
	return categories
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_kubernetes_service_versions":  dataSourceKubernetesServiceVersions(),
		"azurerm_container_group":              dataSourceContainerGroup(),
		"azurerm_container_registry":           dataSourceContainerRegistry(),
		"azurerm_container_registry_token":     dataSourceContainerRegistryToken(),
		"azurerm_container_registry_scope_map": dataSourceContainerRegistryScopeMap(),
		"azurerm_kubernetes_cluster":           dataSourceKubernetesCluster(),
		"azurerm_kubernetes_cluster_node_pool": dataSourceKubernetesClusterNodePool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_container_group":               resourceContainerGroup(),
		"azurerm_container_registry_agent_pool": resourceContainerRegistryAgentPool(),
		"azurerm_container_registry_webhook":    resourceContainerRegistryWebhook(),
		"azurerm_container_registry":            resourceContainerRegistry(),
		"azurerm_container_registry_token":      resourceContainerRegistryToken(),
		"azurerm_container_registry_scope_map":  resourceContainerRegistryScopeMap(),
		"azurerm_kubernetes_cluster":            resourceKubernetesCluster(),
		"azurerm_kubernetes_cluster_node_pool":  resourceKubernetesClusterNodePool(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	dataSources := []sdk.DataSource{
		KubernetesNodePoolSnapshotDataSource{},
	}
	dataSources = append(dataSources, r.autoRegistration.DataSources()...)
	return dataSources
}

func (r Registration) Resources() []sdk.Resource {
	resources := []sdk.Resource{
		ContainerRegistryTaskResource{},
		ContainerRegistryTaskScheduleResource{},
		ContainerRegistryTokenPasswordResource{},
		ContainerConnectedRegistryResource{},
		KubernetesClusterExtensionResource{},
		KubernetesFluxConfigurationResource{},
	}
	resources = append(resources, r.autoRegistration.Resources()...)
	return resources
}
