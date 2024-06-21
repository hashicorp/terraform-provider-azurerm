// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iottimeseriesinsights

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/iot-time-series"
}

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	if !features.FourPointOhBeta() {
		return map[string]*pluginsdk.Resource{
			"azurerm_iot_time_series_insights_access_policy":         resourceIoTTimeSeriesInsightsAccessPolicy(),
			"azurerm_iot_time_series_insights_event_source_eventhub": resourceIoTTimeSeriesInsightsEventSourceEventhub(),
			"azurerm_iot_time_series_insights_event_source_iothub":   resourceIoTTimeSeriesInsightsEventSourceIoTHub(),
			"azurerm_iot_time_series_insights_standard_environment":  resourceIoTTimeSeriesInsightsStandardEnvironment(),
			"azurerm_iot_time_series_insights_gen2_environment":      resourceIoTTimeSeriesInsightsGen2Environment(),
			"azurerm_iot_time_series_insights_reference_data_set":    resourceIoTTimeSeriesInsightsReferenceDataSet(),
		}
	}

	return map[string]*pluginsdk.Resource{}
}
