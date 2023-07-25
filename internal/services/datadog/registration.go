// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datadog

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

type Registration struct{}

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
		"azurerm_datadog_monitor":                   resourceDatadogMonitor(),
		"azurerm_datadog_monitor_tag_rule":          resourceDatadogTagRules(),
		"azurerm_datadog_monitor_sso_configuration": resourceDatadogSingleSignOnConfigurations(),
	}
}
