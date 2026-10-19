// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_analytics_item"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_api_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_smart_detection_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_standard_web_test"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_web_test"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_workbook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/resources/application_insights_workbook_template"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
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
		"azurerm_application_insights": application_insights.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_insights":                      application_insights.RegisterResource(),
		"azurerm_application_insights_analytics_item":       application_insights_analytics_item.RegisterResource(),
		"azurerm_application_insights_api_key":              application_insights_api_key.RegisterResource(),
		"azurerm_application_insights_smart_detection_rule": application_insights_smart_detection_rule.RegisterResource(),
		"azurerm_application_insights_web_test":             application_insights_web_test.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		application_insights_standard_web_test.Resource{},
		application_insights_workbook.Resource{},
		application_insights_workbook_template.Resource{},
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
	return []sdk.FrameworkListWrappedResource{
		application_insights_standard_web_test.ListResource{},
		application_insights_workbook.ListResource{},
		application_insights_workbook_template.ListResource{},
	}
}
