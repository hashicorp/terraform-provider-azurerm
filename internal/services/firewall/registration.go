// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_application_rule_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_nat_rule_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_network_rule_collection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/resources/firewall_policy_rule_collection_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/firewall"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Firewall"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_firewall":        firewall.RegisterDataSource(),
		"azurerm_firewall_policy": firewall_policy.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_firewall": firewall.RegisterResource(),
		"azurerm_firewall_application_rule_collection":  firewall_application_rule_collection.RegisterResource(),
		"azurerm_firewall_nat_rule_collection":          firewall_nat_rule_collection.RegisterResource(),
		"azurerm_firewall_network_rule_collection":      firewall_network_rule_collection.RegisterResource(),
		"azurerm_firewall_policy":                       firewall_policy.RegisterResource(),
		"azurerm_firewall_policy_rule_collection_group": firewall_policy_rule_collection_group.RegisterResource(),
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
		firewall.ListResource{},
		firewall_policy.ListResource{},
		firewall_policy_rule_collection_group.ListResource{},
	}
}
