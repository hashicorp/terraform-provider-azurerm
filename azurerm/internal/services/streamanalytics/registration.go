package streamanalytics

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Stream Analytics"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Stream Analytics",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_stream_analytics_job": dataSourceArmStreamAnalyticsJob(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_stream_analytics_job":                     resourceStreamAnalyticsJob(),
		"azurerm_stream_analytics_function_javascript_udf": resourceStreamAnalyticsFunctionUDF(),
		"azurerm_stream_analytics_output_blob":             resourceStreamAnalyticsOutputBlob(),
		"azurerm_stream_analytics_output_mssql":            resourceStreamAnalyticsOutputSql(),
		"azurerm_stream_analytics_output_eventhub":         resourceStreamAnalyticsOutputEventHub(),
		"azurerm_stream_analytics_output_servicebus_queue": resourceStreamAnalyticsOutputServiceBusQueue(),
		"azurerm_stream_analytics_output_servicebus_topic": resourceStreamAnalyticsOutputServiceBusTopic(),
		"azurerm_stream_analytics_reference_input_blob":    resourceStreamAnalyticsReferenceInputBlob(),
		"azurerm_stream_analytics_stream_input_blob":       resourceStreamAnalyticsStreamInputBlob(),
		"azurerm_stream_analytics_stream_input_eventhub":   resourceStreamAnalyticsStreamInputEventHub(),
		"azurerm_stream_analytics_stream_input_iothub":     resourceStreamAnalyticsStreamInputIoTHub(),
	}
}
