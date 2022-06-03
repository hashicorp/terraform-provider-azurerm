package cdn

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
		"azurerm_cdn_profile": dataSourceCdnProfile(),

		// FrontDoor
		"azurerm_cdn_frontdoor_endpoint": dataSourceCdnFrontDoorEndpoint(),
		"azurerm_cdn_frontdoor_profile":  dataSourceCdnFrontDoorProfile(),
		"azurerm_cdn_frontdoor_rule_set": dataSourceCdnFrontDoorRuleSet(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		// CDN
		"azurerm_cdn_endpoint":               resourceCdnEndpoint(),
		"azurerm_cdn_endpoint_custom_domain": resourceArmCdnEndpointCustomDomain(),
		"azurerm_cdn_profile":                resourceCdnProfile(),

		// FrontDoor
		"azurerm_cdn_frontdoor_endpoint": resourceCdnFrontDoorEndpoint(),
		"azurerm_cdn_frontdoor_profile":  resourceCdnFrontDoorProfile(),
		"azurerm_cdn_frontdoor_rule_set": resourceCdnFrontDoorRuleSet(),
	}
}
