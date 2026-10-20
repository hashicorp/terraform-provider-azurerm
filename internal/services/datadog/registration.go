// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/resources/datadog_monitor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/resources/datadog_monitor_sso_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/resources/datadog_monitor_tag_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.FrameworkServiceRegistration = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/datadog"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Datadog"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Datadog",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_datadog_monitor":                   datadog_monitor.RegisterResource(),
		"azurerm_datadog_monitor_sso_configuration": datadog_monitor_sso_configuration.RegisterResource(),
		"azurerm_datadog_monitor_tag_rule":          datadog_monitor_tag_rule.RegisterResource(),
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
	return []sdk.FrameworkListWrappedResource{}
}
