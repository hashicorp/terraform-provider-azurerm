// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_connected_registry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_agent_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_cache_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_credential_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_scope_map"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_task"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_task_schedule_run_now"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_token"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_token_password"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/container_registry_webhook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_automatic_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster_deployment_safeguard"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster_extension"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_cluster_node_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_fleet_manager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_fleet_update_run"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_fleet_update_strategy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_flux_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_node_pool_snapshot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/resources/kubernetes_service_versions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct {
	autoRegistration
}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
	_ sdk.UntypedServiceRegistration   = Registration{}
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
		"azurerm_container_group":              container_group.RegisterDataSource(),
		"azurerm_container_registry":           container_registry.RegisterDataSource(),
		"azurerm_container_registry_scope_map": container_registry_scope_map.RegisterDataSource(),
		"azurerm_container_registry_token":     container_registry_token.RegisterDataSource(),
		"azurerm_kubernetes_cluster":           kubernetes_cluster.RegisterDataSource(),
		"azurerm_kubernetes_cluster_node_pool": kubernetes_cluster_node_pool.RegisterDataSource(),
		"azurerm_kubernetes_service_versions":  kubernetes_service_versions.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_container_group":               container_group.RegisterResource(),
		"azurerm_container_registry":            container_registry.RegisterResource(),
		"azurerm_container_registry_agent_pool": container_registry_agent_pool.RegisterResource(),
		"azurerm_container_registry_scope_map":  container_registry_scope_map.RegisterResource(),
		"azurerm_container_registry_token":      container_registry_token.RegisterResource(),
		"azurerm_container_registry_webhook":    container_registry_webhook.RegisterResource(),
		"azurerm_kubernetes_cluster":            kubernetes_cluster.RegisterResource(),
		"azurerm_kubernetes_cluster_node_pool":  kubernetes_cluster_node_pool.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	dataSources := []sdk.DataSource{
		container_registry_cache_rule.DataSource{},
		kubernetes_automatic_cluster.DataSource{},
		kubernetes_fleet_manager.DataSource{},
		kubernetes_node_pool_snapshot.DataSource{},
	}
	dataSources = append(dataSources, r.autoRegistration.DataSources()...)
	return dataSources
}

func (r Registration) Resources() []sdk.Resource {
	resources := []sdk.Resource{
		container_connected_registry.Resource{},
		container_registry_cache_rule.Resource{},
		container_registry_credential_set.Resource{},
		container_registry_task.Resource{},
		container_registry_task_schedule_run_now.Resource{},
		container_registry_token_password.Resource{},
		kubernetes_automatic_cluster.Resource{},
		kubernetes_cluster_deployment_safeguard.Resource{},
		kubernetes_cluster_extension.Resource{},
		kubernetes_fleet_manager.Resource{},
		kubernetes_fleet_update_run.Resource{},
		kubernetes_fleet_update_strategy.Resource{},
		kubernetes_flux_configuration.Resource{},
	}
	resources = append(resources, r.autoRegistration.Resources()...)
	return resources
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
	return []sdk.FrameworkListWrappedResource{
		kubernetes_automatic_cluster.ListResource{},
		kubernetes_cluster_node_pool.ListResource{},
	}
}
