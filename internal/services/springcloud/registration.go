// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_accelerator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_active_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_api_portal"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_api_portal_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_app_cosmosdb_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_app_dynamics_application_performance_monitoring"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_app_mysql_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_app_redis_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_application_insights_application_performance_monitoring"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_application_live_view"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_build_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_build_pack_binding"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_builder"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_configuration_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_container_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_customized_accelerator"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_dev_tool_portal"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_dynatrace_application_performance_monitoring"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_elastic_application_performance_monitoring"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_gateway_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_gateway_route_config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_java_deployment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_new_relic_application_performance_monitoring"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/resources/spring_cloud_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/spring"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Spring Cloud"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Spring Cloud",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_spring_cloud_app":     spring_cloud_app.RegisterDataSource(),
		"azurerm_spring_cloud_service": spring_cloud_service.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_spring_cloud_active_deployment":        spring_cloud_active_deployment.RegisterResource(),
		"azurerm_spring_cloud_api_portal_custom_domain": spring_cloud_api_portal_custom_domain.RegisterResource(),
		"azurerm_spring_cloud_app":                      spring_cloud_app.RegisterResource(),
		"azurerm_spring_cloud_app_cosmosdb_association": spring_cloud_app_cosmosdb_association.RegisterResource(),
		"azurerm_spring_cloud_app_mysql_association":    spring_cloud_app_mysql_association.RegisterResource(),
		"azurerm_spring_cloud_app_redis_association":    spring_cloud_app_redis_association.RegisterResource(),
		"azurerm_spring_cloud_build_deployment":         spring_cloud_build_deployment.RegisterResource(),
		"azurerm_spring_cloud_build_pack_binding":       spring_cloud_build_pack_binding.RegisterResource(),
		"azurerm_spring_cloud_builder":                  spring_cloud_builder.RegisterResource(),
		"azurerm_spring_cloud_certificate":              spring_cloud_certificate.RegisterResource(),
		"azurerm_spring_cloud_container_deployment":     spring_cloud_container_deployment.RegisterResource(),
		"azurerm_spring_cloud_custom_domain":            spring_cloud_custom_domain.RegisterResource(),
		"azurerm_spring_cloud_gateway_custom_domain":    spring_cloud_gateway_custom_domain.RegisterResource(),
		"azurerm_spring_cloud_gateway_route_config":     spring_cloud_gateway_route_config.RegisterResource(),
		"azurerm_spring_cloud_java_deployment":          spring_cloud_java_deployment.RegisterResource(),
		"azurerm_spring_cloud_service":                  spring_cloud_service.RegisterResource(),
		"azurerm_spring_cloud_storage":                  spring_cloud_storage.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		spring_cloud_accelerator.Resource{},
		spring_cloud_api_portal.Resource{},
		spring_cloud_app_dynamics_application_performance_monitoring.Resource{},
		spring_cloud_application_insights_application_performance_monitoring.Resource{},
		spring_cloud_application_live_view.Resource{},
		spring_cloud_configuration_service.Resource{},
		spring_cloud_customized_accelerator.Resource{},
		spring_cloud_dev_tool_portal.Resource{},
		spring_cloud_dynatrace_application_performance_monitoring.Resource{},
		spring_cloud_elastic_application_performance_monitoring.Resource{},
		spring_cloud_gateway.Resource{},
		spring_cloud_new_relic_application_performance_monitoring.Resource{},
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
