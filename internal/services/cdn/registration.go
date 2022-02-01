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
		"azurerm_frontdoor_profile_origin_group":    resourceFrontdoorProfileOriginGroup(),
		"azurerm_frontdoor_profile_custom_domain":   resourceFrontdoorProfileCustomDomain(),
		"azurerm_frontdoor_profile_origin":          resourceFrontdoorProfileOrigin(),
		"azurerm_frontdoor_profile_endpoint":        resourceFrontdoorProfileEndpoint(),
		"azurerm_frontdoor_profile_policy":          resourceFrontdoorProfilePolicy(),
		"azurerm_frontdoor_profile_rule_set":        resourceFrontdoorProfileRuleSet(),
		"azurerm_frontdoor_profile_security_policy": resourceFrontdoorProfileSecurityPolicy(),
		"azurerm_frontdoor_profile":                 resourceFrontdoorProfile(),
		"azurerm_frontdoor_profile_route":           resourceFrontdoorProfileRoute(),
		"azurerm_frontdoor_profile_secret":          resourceFrontdoorProfileSecret(),
		"azurerm_frontdoor_profile_rule":            resourceFrontdoorProfileRule(),
		"azurerm_cdn_endpoint":                      resourceCdnEndpoint(),
		"azurerm_cdn_endpoint_custom_domain":        resourceArmCdnEndpointCustomDomain(),
		"azurerm_cdn_profile":                       resourceCdnProfile(),
	}
}
