// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/application-insights"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Application Insights"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Application Insights",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_insights": dataSourceApplicationInsights(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_insights_api_key":              resourceApplicationInsightsAPIKey(),
		"azurerm_application_insights":                      resourceApplicationInsights(),
		"azurerm_application_insights_analytics_item":       resourceApplicationInsightsAnalyticsItem(),
		"azurerm_application_insights_smart_detection_rule": resourceApplicationInsightsSmartDetectionRule(),

		// TODO change in 4.0 to azurerm_application_insights_classic_web_test
		"azurerm_application_insights_web_test": resourceApplicationInsightsWebTests(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ApplicationInsightsWorkbookResource{},
		ApplicationInsightsWorkbookTemplateResource{},
		ApplicationInsightsStandardWebTestResource{},
	}
}
