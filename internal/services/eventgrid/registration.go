// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_domain"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_domain_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_event_subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_namespace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_namespace_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_partner_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_partner_namespace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_partner_registration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_system_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_system_topic_event_subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/resources/eventgrid_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
	return []sdk.DataSource{
		eventgrid_partner_namespace.DataSource{},
		eventgrid_partner_registration.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		eventgrid_namespace.Resource{},
		eventgrid_namespace_topic.Resource{},
		eventgrid_partner_configuration.Resource{},
		eventgrid_partner_namespace.Resource{},
		eventgrid_partner_registration.Resource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventgrid_domain":       eventgrid_domain.RegisterDataSource(),
		"azurerm_eventgrid_domain_topic": eventgrid_domain_topic.RegisterDataSource(),
		"azurerm_eventgrid_system_topic": eventgrid_system_topic.RegisterDataSource(),
		"azurerm_eventgrid_topic":        eventgrid_topic.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventgrid_domain":                          eventgrid_domain.RegisterResource(),
		"azurerm_eventgrid_domain_topic":                    eventgrid_domain_topic.RegisterResource(),
		"azurerm_eventgrid_event_subscription":              eventgrid_event_subscription.RegisterResource(),
		"azurerm_eventgrid_system_topic":                    eventgrid_system_topic.RegisterResource(),
		"azurerm_eventgrid_system_topic_event_subscription": eventgrid_system_topic_event_subscription.RegisterResource(),
		"azurerm_eventgrid_topic":                           eventgrid_topic.RegisterResource(),
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
	return []sdk.FrameworkListWrappedResource{}
}
