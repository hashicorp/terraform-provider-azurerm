package applicationinsights

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_application_insights_web_test":             resourceApplicationInsightsWebTests(),
	}
}
