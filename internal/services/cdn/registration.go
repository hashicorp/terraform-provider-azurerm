package cdn

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

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
		"azurerm_cdn_profile": dataSourceCdnProfile(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_frontdoor_origin_group":     resourceFrontdoorOriginGroup(),
		"azurerm_frontdoor_custom_domain":    resourceFrontdoorCustomDomain(),
		"azurerm_frontdoor_origin":           resourceFrontdoorOrigin(),
		"azurerm_frontdoor_endpoint":         resourceFrontdoorEndpoint(),
		"azurerm_frontdoor_rule_set":         resourceFrontdoorRuleSet(),
		"azurerm_frontdoor_security_policy":  resourceFrontdoorSecurityPolicy(),
		"azurerm_frontdoor_profile":          resourceFrontdoorProfile(),
		"azurerm_frontdoor_route":            resourceFrontdoorRoute(),
		"azurerm_frontdoor_secret":           resourceFrontdoorSecret(),
		"azurerm_frontdoor_rule":             resourceFrontdoorRule(),
		"azurerm_cdn_endpoint":               resourceCdnEndpoint(),
		"azurerm_cdn_endpoint_custom_domain": resourceArmCdnEndpointCustomDomain(),
		"azurerm_cdn_profile":                resourceCdnProfile(),
	}
}
