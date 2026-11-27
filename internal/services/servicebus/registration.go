// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_namespace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_namespace_authorization_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_namespace_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_namespace_disaster_recovery_config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_queue"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_queue_authorization_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_subscription_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_topic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/resources/servicebus_topic_authorization_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/service-bus"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "ServiceBus"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		// TODO: change this to ServiceBus
		"Messaging",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_servicebus_namespace":                          servicebus_namespace.RegisterDataSource(),
		"azurerm_servicebus_namespace_authorization_rule":       servicebus_namespace_authorization_rule.RegisterDataSource(),
		"azurerm_servicebus_namespace_disaster_recovery_config": servicebus_namespace_disaster_recovery_config.RegisterDataSource(),
		"azurerm_servicebus_queue":                              servicebus_queue.RegisterDataSource(),
		"azurerm_servicebus_queue_authorization_rule":           servicebus_queue_authorization_rule.RegisterDataSource(),
		"azurerm_servicebus_subscription":                       servicebus_subscription.RegisterDataSource(),
		"azurerm_servicebus_topic":                              servicebus_topic.RegisterDataSource(),
		"azurerm_servicebus_topic_authorization_rule":           servicebus_topic_authorization_rule.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_servicebus_namespace":                          servicebus_namespace.RegisterResource(),
		"azurerm_servicebus_namespace_authorization_rule":       servicebus_namespace_authorization_rule.RegisterResource(),
		"azurerm_servicebus_namespace_disaster_recovery_config": servicebus_namespace_disaster_recovery_config.RegisterResource(),
		"azurerm_servicebus_queue":                              servicebus_queue.RegisterResource(),
		"azurerm_servicebus_queue_authorization_rule":           servicebus_queue_authorization_rule.RegisterResource(),
		"azurerm_servicebus_subscription":                       servicebus_subscription.RegisterResource(),
		"azurerm_servicebus_subscription_rule":                  servicebus_subscription_rule.RegisterResource(),
		"azurerm_servicebus_topic":                              servicebus_topic.RegisterResource(),
		"azurerm_servicebus_topic_authorization_rule":           servicebus_topic_authorization_rule.RegisterResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		servicebus_namespace_customer_managed_key.Resource{},
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
	return []sdk.FrameworkListWrappedResource{
		servicebus_namespace_customer_managed_key.ListResource{},
		servicebus_namespace.ListResource{},
		servicebus_queue.ListResource{},
	}
}
