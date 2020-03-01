package eventgrid

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "EventGrid"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Messaging",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventgrid_topic": dataSourceArmEventGridTopic(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventgrid_domain":             resourceArmEventGridDomain(),
		"azurerm_eventgrid_domain_topic":       resourceArmEventGridDomainTopic(),
		"azurerm_eventgrid_event_subscription": resourceArmEventGridEventSubscription(),
		"azurerm_eventgrid_topic":              resourceArmEventGridTopic(),
	}
}
