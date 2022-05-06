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
		"azurerm_cdn_profile": dataSourceCdnProfile(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_cdn_frontdoor_origin_group":                   resourceCdnFrontdoorOriginGroup(),
		"azurerm_cdn_frontdoor_custom_domain":                  resourceCdnFrontdoorCustomDomain(),
		"azurerm_cdn_frontdoor_custom_domain_txt_validator":    resourceCdnFrontdoorCustomDomainTxtValidator(),
		"azurerm_cdn_frontdoor_custom_domain_secret_validator": resourceCdnFrontdoorCustomDomainSecretValidator(),
		"azurerm_cdn_frontdoor_origin":                         resourceCdnFrontdoorOrigin(),
		"azurerm_cdn_frontdoor_endpoint":                       resourceCdnFrontdoorEndpoint(),
		"azurerm_cdn_frontdoor_rule_set":                       resourceCdnFrontdoorRuleSet(),
		"azurerm_cdn_frontdoor_security_policy":                resourceCdnFrontdoorSecurityPolicy(),
		"azurerm_cdn_frontdoor_profile":                        resourceCdnFrontdoorProfile(),
		"azurerm_cdn_frontdoor_route":                          resourceCdnFrontdoorRoute(),
		"azurerm_cdn_frontdoor_secret":                         resourceCdnFrontdoorSecret(),
		"azurerm_cdn_frontdoor_rule":                           resourceCdnFrontdoorRule(),
		"azurerm_cdn_frontdoor_firewall_policy":                resourceCdnFrontdoorFirewallPolicy(),
		"azurerm_cdn_endpoint":                                 resourceCdnEndpoint(),
		"azurerm_cdn_endpoint_custom_domain":                   resourceArmCdnEndpointCustomDomain(),
		"azurerm_cdn_profile":                                  resourceCdnProfile(),
	}
}
