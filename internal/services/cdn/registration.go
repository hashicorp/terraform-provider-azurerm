// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/actions/azure_front_door_cache_purge"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_endpoint_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_batch_rule_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_custom_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_custom_domain_association"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_firewall_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_origin"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_origin_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_profile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_route"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_rule_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_secret"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_security_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_profile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

var _ sdk.FrameworkServiceRegistration = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/cdn"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "CDN"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"CDN",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		// CDN
		"azurerm_cdn_profile": cdn_profile.RegisterDataSource(),

		// FrontDoor
		"azurerm_cdn_frontdoor_custom_domain":   cdn_frontdoor_custom_domain.RegisterDataSource(),
		"azurerm_cdn_frontdoor_endpoint":        cdn_frontdoor_endpoint.RegisterDataSource(),
		"azurerm_cdn_frontdoor_firewall_policy": cdn_frontdoor_firewall_policy.RegisterDataSource(),
		"azurerm_cdn_frontdoor_origin_group":    cdn_frontdoor_origin_group.RegisterDataSource(),
		"azurerm_cdn_frontdoor_profile":         cdn_frontdoor_profile.RegisterDataSource(),
		"azurerm_cdn_frontdoor_rule_set":        cdn_frontdoor_rule_set.RegisterDataSource(),
		"azurerm_cdn_frontdoor_secret":          cdn_frontdoor_secret.RegisterDataSource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		cdn_frontdoor_batch_rule_set.DataSource{},
		cdn_frontdoor_security_policy.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		cdn_frontdoor_batch_rule_set.Resource{},
		cdn_frontdoor_rule.Resource{},
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		// CDN
		"azurerm_cdn_endpoint":               cdn_endpoint.RegisterResource(),
		"azurerm_cdn_endpoint_custom_domain": cdn_endpoint_custom_domain.RegisterResource(),
		"azurerm_cdn_profile":                cdn_profile.RegisterResource(),

		// FrontDoor
		"azurerm_cdn_frontdoor_custom_domain":             cdn_frontdoor_custom_domain.RegisterResource(),
		"azurerm_cdn_frontdoor_custom_domain_association": cdn_frontdoor_custom_domain_association.RegisterResource(),
		"azurerm_cdn_frontdoor_endpoint":                  cdn_frontdoor_endpoint.RegisterResource(),
		"azurerm_cdn_frontdoor_firewall_policy":           cdn_frontdoor_firewall_policy.RegisterResource(),
		"azurerm_cdn_frontdoor_origin":                    cdn_frontdoor_origin.RegisterResource(),
		"azurerm_cdn_frontdoor_origin_group":              cdn_frontdoor_origin_group.RegisterResource(),
		"azurerm_cdn_frontdoor_profile":                   cdn_frontdoor_profile.RegisterResource(),
		"azurerm_cdn_frontdoor_route":                     cdn_frontdoor_route.RegisterResource(),
		"azurerm_cdn_frontdoor_rule_set":                  cdn_frontdoor_rule_set.RegisterResource(),
		"azurerm_cdn_frontdoor_secret":                    cdn_frontdoor_secret.RegisterResource(),
		"azurerm_cdn_frontdoor_security_policy":           cdn_frontdoor_security_policy.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		azure_front_door_cache_purge.RegisterAction,
	}
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
		cdn_frontdoor_batch_rule_set.ListResource{},
		cdn_frontdoor_origin_group.ListResource{},
		cdn_frontdoor_origin.ListResource{},
	}
}
