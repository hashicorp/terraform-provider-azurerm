// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/actions/web_app_set_slot_distribution"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/app_service_environment_v3"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/app_service_source_control"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/app_service_source_control_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/function_app_active_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/function_app_flex_consumption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/function_app_function"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/function_app_hybrid_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_function_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_function_app_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_web_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_web_app_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/service_plan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/source_control_token"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/static_web_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/static_web_app_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/static_web_app_function_app_registration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/web_app_active_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/web_app_hybrid_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_function_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_function_app_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_web_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_web_app_slot"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/app-service"
}

func (r Registration) WebsiteCategories() []string {
	return nil
}

func (r Registration) Name() string {
	return "AppService"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		app_service_environment_v3.DataSource{},
		source_control_token.DataSource{},
		linux_function_app.DataSource{},
		linux_web_app.DataSource{},
		service_plan.DataSource{},
		static_web_app.DataSource{},
		windows_function_app.DataSource{},
		windows_web_app.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		app_service_environment_v3.Resource{},
		source_control_token.Resource{},
		function_app_active_slot.Resource{},
		function_app_flex_consumption.Resource{},
		function_app_function.Resource{},
		function_app_hybrid_connection.Resource{},
		linux_function_app.Resource{},
		linux_function_app_slot.Resource{},
		linux_web_app.Resource{},
		linux_web_app_slot.Resource{},
		service_plan.Resource{},
		app_service_source_control.Resource{},
		app_service_source_control_slot.Resource{},
		static_web_app_custom_domain.Resource{},
		static_web_app_function_app_registration.Resource{},
		static_web_app.Resource{},
		web_app_active_slot.Resource{},
		web_app_hybrid_connection.Resource{},
		windows_function_app.Resource{},
		windows_function_app_slot.Resource{},
		windows_web_app.Resource{},
		windows_web_app_slot.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		web_app_set_slot_distribution.RegisterAction,
	}
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
		&service_plan.ListResource{},
	}
}
