// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_certificate_binding"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_certificate_order"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_custom_hostname_binding"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_managed_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_public_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_slot_custom_hostname_binding"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_slot_virtual_network_swift_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/app_service_virtual_network_swift_connection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/resources/function_app_host_keys"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
	_ sdk.UntypedServiceRegistration   = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "Web"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"App Service (Web Apps)",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_app_service_certificate":       app_service_certificate.RegisterDataSource(),
		"azurerm_app_service_certificate_order": app_service_certificate_order.RegisterDataSource(),
		"azurerm_function_app_host_keys":        function_app_host_keys.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_app_service_certificate":                           app_service_certificate.RegisterResource(),
		"azurerm_app_service_certificate_binding":                   app_service_certificate_binding.RegisterResource(),
		"azurerm_app_service_certificate_order":                     app_service_certificate_order.RegisterResource(),
		"azurerm_app_service_custom_hostname_binding":               app_service_custom_hostname_binding.RegisterResource(),
		"azurerm_app_service_managed_certificate":                   app_service_managed_certificate.RegisterResource(),
		"azurerm_app_service_public_certificate":                    app_service_public_certificate.RegisterResource(),
		"azurerm_app_service_slot_custom_hostname_binding":          app_service_slot_custom_hostname_binding.RegisterResource(),
		"azurerm_app_service_slot_virtual_network_swift_connection": app_service_slot_virtual_network_swift_connection.RegisterResource(),
		"azurerm_app_service_virtual_network_swift_connection":      app_service_virtual_network_swift_connection.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
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
		app_service_custom_hostname_binding.ListResource{},
	}
}
