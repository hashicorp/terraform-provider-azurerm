// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/frontdoor"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "FrontDoor"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_frontdoor":                            resourceFrontDoor(),
		"azurerm_frontdoor_firewall_policy":            resourceFrontDoorFirewallPolicy(),
		"azurerm_frontdoor_custom_https_configuration": resourceFrontDoorCustomHTTPSConfiguration(),
		"azurerm_frontdoor_rules_engine":               resourceFrontDoorRulesEngine(),
	}
}
