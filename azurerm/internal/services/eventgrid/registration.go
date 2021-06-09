package eventgrid

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventgrid_topic":        dataSourceEventGridTopic(),
		"azurerm_eventgrid_domain_topic": dataSourceEventGridDomainTopic(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventgrid_domain":                          resourceEventGridDomain(),
		"azurerm_eventgrid_domain_topic":                    resourceEventGridDomainTopic(),
		"azurerm_eventgrid_event_subscription":              resourceEventGridEventSubscription(),
		"azurerm_eventgrid_topic":                           resourceEventGridTopic(),
		"azurerm_eventgrid_system_topic":                    resourceEventGridSystemTopic(),
		"azurerm_eventgrid_system_topic_event_subscription": resourceEventGridSystemTopicEventSubscription(),
	}
}
