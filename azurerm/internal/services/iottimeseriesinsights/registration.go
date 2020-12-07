package iottimeseriesinsights

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Time Series Insights"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Time Series Insights",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_iot_time_series_insights_access_policy":        resourceArmIoTTimeSeriesInsightsAccessPolicy(),
		"azurerm_iot_time_series_insights_standard_environment": resourceArmIoTTimeSeriesInsightsStandardEnvironment(),
		"azurerm_iot_time_series_insights_gen2_environment":     resourceArmIoTTimeSeriesInsightsGen2Environment(),
		"azurerm_iot_time_series_insights_reference_data_set":   resourceArmIoTTimeSeriesInsightsReferenceDataSet(),
	}
}
