// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/iot-hub"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "IoT Hub"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"IoT Hub",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_iothub_dps":                      dataSourceIotHubDPS(),
		"azurerm_iothub_dps_shared_access_policy": dataSourceIotHubDPSSharedAccessPolicy(),
		"azurerm_iothub_shared_access_policy":     dataSourceIotHubSharedAccessPolicy(),
		"azurerm_iothub":                          dataSourceIotHub(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_iothub_certificate":                resourceIotHubCertificate(),
		"azurerm_iothub_dps":                        resourceIotHubDPS(),
		"azurerm_iothub_dps_certificate":            resourceIotHubDPSCertificate(),
		"azurerm_iothub_dps_shared_access_policy":   resourceIotHubDPSSharedAccessPolicy(),
		"azurerm_iothub_consumer_group":             resourceIotHubConsumerGroup(),
		"azurerm_iothub":                            resourceIotHub(),
		"azurerm_iothub_fallback_route":             resourceIotHubFallbackRoute(),
		"azurerm_iothub_enrichment":                 resourceIotHubEnrichment(),
		"azurerm_iothub_route":                      resourceIotHubRoute(),
		"azurerm_iothub_endpoint_eventhub":          resourceIotHubEndpointEventHub(),
		"azurerm_iothub_endpoint_servicebus_queue":  resourceIotHubEndpointServiceBusQueue(),
		"azurerm_iothub_endpoint_servicebus_topic":  resourceIotHubEndpointServiceBusTopic(),
		"azurerm_iothub_endpoint_storage_container": resourceIotHubEndpointStorageContainer(),
		"azurerm_iothub_shared_access_policy":       resourceIotHubSharedAccessPolicy(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		IotHubDeviceUpdateAccountResource{},
		IotHubDeviceUpdateInstanceResource{},
		IotHubFileUploadResource{},
		IotHubEndpointCosmosDBAccountResource{},
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
