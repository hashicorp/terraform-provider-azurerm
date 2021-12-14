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
		"azurerm_cdn_endpoint":                 resourceCdnEndpoint(),
		"azurerm_cdn_endpoint_custom_domain":   resourceArmCdnEndpointCustomDomain(),
		"azurerm_cdn_profile":                  resourceCdnProfile(),
		"azurerm_cdn_frontdoor_profile":        resourceCdnProfile(), // re-use of azurerm_cdn_profile
		"azurerm_cdn_frontdoor_endpoint":       resourceAfdEndpoints(),
		"azurerm_cdn_frontdoor_endpoint_route": resourceAfdEndpointRoutes(),
		"azurerm_cdn_frontdoor_origin_group":   resourceAfdOriginGroups(),
		"azurerm_cdn_frontdoor_origin":         resourceAfdOrigin(),
		"azurerm_cdn_frontdoor_custom_domain":  resourceAfdCustomDomains(),
	}
}
