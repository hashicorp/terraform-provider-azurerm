// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/event-grid"
}

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

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		EventGridNamespaceResource{},
		EventGridPartnerConfigurationResource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventgrid_topic":        dataSourceEventGridTopic(),
		"azurerm_eventgrid_domain":       dataSourceEventGridDomain(),
		"azurerm_eventgrid_domain_topic": dataSourceEventGridDomainTopic(),
		"azurerm_eventgrid_system_topic": dataSourceEventGridSystemTopic(),
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
