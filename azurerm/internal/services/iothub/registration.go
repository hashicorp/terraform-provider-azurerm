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
		"azurerm_iothub_dps":                      dataSourceArmIotHubDPS(),
		"azurerm_iothub_dps_shared_access_policy": dataSourceIotHubDPSSharedAccessPolicy(),
		"azurerm_iothub_shared_access_policy":     dataSourceArmIotHubSharedAccessPolicy(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_iothub_dps":                        resourceArmIotHubDPS(),
		"azurerm_iothub_dps_certificate":            resourceArmIotHubDPSCertificate(),
		"azurerm_iothub_dps_shared_access_policy":   resourceArmIotHubDPSSharedAccessPolicy(),
		"azurerm_iothub_consumer_group":             resourceArmIotHubConsumerGroup(),
		"azurerm_iothub":                            resourceArmIotHub(),
		"azurerm_iothub_fallback_route":             resourceArmIotHubFallbackRoute(),
		"azurerm_iothub_route":                      resourceArmIotHubRoute(),
		"azurerm_iothub_endpoint_eventhub":          resourceArmIotHubEndpointEventHub(),
		"azurerm_iothub_endpoint_servicebus_queue":  resourceArmIotHubEndpointServiceBusQueue(),
		"azurerm_iothub_endpoint_servicebus_topic":  resourceArmIotHubEndpointServiceBusTopic(),
		"azurerm_iothub_endpoint_storage_container": resourceArmIotHubEndpointStorageContainer(),
		"azurerm_iothub_shared_access_policy":       resourceArmIotHubSharedAccessPolicy()}
}
