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
		"azurerm_stream_analytics_job":                     resourceArmStreamAnalyticsJob(),
		"azurerm_stream_analytics_function_javascript_udf": resourceArmStreamAnalyticsFunctionUDF(),
		"azurerm_stream_analytics_output_blob":             resourceArmStreamAnalyticsOutputBlob(),
		"azurerm_stream_analytics_output_mssql":            resourceArmStreamAnalyticsOutputSql(),
		"azurerm_stream_analytics_output_eventhub":         resourceArmStreamAnalyticsOutputEventHub(),
		"azurerm_stream_analytics_output_servicebus_queue": resourceArmStreamAnalyticsOutputServiceBusQueue(),
		"azurerm_stream_analytics_output_servicebus_topic": resourceArmStreamAnalyticsOutputServiceBusTopic(),
		"azurerm_stream_analytics_reference_input_blob":    resourceArmStreamAnalyticsReferenceInputBlob(),
		"azurerm_stream_analytics_stream_input_blob":       resourceArmStreamAnalyticsStreamInputBlob(),
		"azurerm_stream_analytics_stream_input_eventhub":   resourceArmStreamAnalyticsStreamInputEventHub(),
		"azurerm_stream_analytics_stream_input_iothub":     resourceArmStreamAnalyticsStreamInputIoTHub(),
	}
}
