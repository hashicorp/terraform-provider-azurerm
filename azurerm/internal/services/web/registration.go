package web

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_app_service":                   dataSourceAppService(),
		"azurerm_app_service_certificate_order": dataSourceAppServiceCertificateOrder(),
		"azurerm_app_service_environment":       dataSourceAppServiceEnvironment(),
		"azurerm_app_service_certificate":       dataSourceAppServiceCertificate(),
		"azurerm_app_service_plan":              dataSourceAppServicePlan(),
		"azurerm_function_app":                  dataSourceFunctionApp(),
		"azurerm_function_app_host_keys":        dataSourceFunctionAppHostKeys(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_app_service_active_slot":                           resourceAppServiceActiveSlot(),
		"azurerm_app_service_certificate":                           resourceAppServiceCertificate(),
		"azurerm_app_service_certificate_order":                     resourceAppServiceCertificateOrder(),
		"azurerm_app_service_custom_hostname_binding":               resourceAppServiceCustomHostnameBinding(),
		"azurerm_app_service_certificate_binding":                   resourceAppServiceCertificateBinding(),
		"azurerm_app_service_environment":                           resourceAppServiceEnvironment(),
		"azurerm_app_service_hybrid_connection":                     resourceAppServiceHybridConnection(),
		"azurerm_app_service_managed_certificate":                   resourceAppServiceManagedCertificate(),
		"azurerm_app_service_plan":                                  resourceAppServicePlan(),
		"azurerm_app_service_slot":                                  resourceAppServiceSlot(),
		"azurerm_app_service_slot_virtual_network_swift_connection": resourceAppServiceSlotVirtualNetworkSwiftConnection(),
		"azurerm_app_service_source_control_token":                  resourceAppServiceSourceControlToken(),
		"azurerm_app_service_virtual_network_swift_connection":      resourceAppServiceVirtualNetworkSwiftConnection(),
		"azurerm_app_service":                                       resourceAppService(),
		"azurerm_function_app":                                      resourceFunctionApp(),
		"azurerm_function_app_slot":                                 resourceFunctionAppSlot(),
	}
}
