package streamanalytics

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}
var _ sdk.UntypedServiceRegistration = Registration{}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		OutputTableResource{},
		ClusterResource{},
		ManagedPrivateEndpointResource{},
	}
}

func (r Registration) Name() string {
	return "Stream Analytics"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Stream Analytics",
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_stream_analytics_job": dataSourceArmStreamAnalyticsJob(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_stream_analytics_job":                     resourceStreamAnalyticsJob(),
		"azurerm_stream_analytics_function_javascript_udf": resourceStreamAnalyticsFunctionUDF(),
		"azurerm_stream_analytics_output_blob":             resourceStreamAnalyticsOutputBlob(),
		"azurerm_stream_analytics_output_mssql":            resourceStreamAnalyticsOutputSql(),
		"azurerm_stream_analytics_output_eventhub":         resourceStreamAnalyticsOutputEventHub(),
		"azurerm_stream_analytics_output_servicebus_queue": resourceStreamAnalyticsOutputServiceBusQueue(),
		"azurerm_stream_analytics_output_servicebus_topic": resourceStreamAnalyticsOutputServiceBusTopic(),
		"azurerm_stream_analytics_reference_input_blob":    resourceStreamAnalyticsReferenceInputBlob(),
		"azurerm_stream_analytics_reference_input_mssql":   resourceStreamAnalyticsReferenceMsSql(),
		"azurerm_stream_analytics_stream_input_blob":       resourceStreamAnalyticsStreamInputBlob(),
		"azurerm_stream_analytics_stream_input_eventhub":   resourceStreamAnalyticsStreamInputEventHub(),
		"azurerm_stream_analytics_stream_input_iothub":     resourceStreamAnalyticsStreamInputIoTHub(),
	}
}
