package iothub

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_iothub_dps":                      dataSourceIotHubDPS(),
		"azurerm_iothub_dps_shared_access_policy": dataSourceIotHubDPSSharedAccessPolicy(),
		"azurerm_iothub_shared_access_policy":     dataSourceIotHubSharedAccessPolicy(),
		"azurerm_iothub":                          dataSourceIotHub(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
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
