// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_authorization_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_consumer_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_namespace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_namespace_authorization_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_namespace_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_namespace_disaster_recovery_config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_namespace_schema_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/resources/eventhub_sas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/event-hubs"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "EventHub"
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
		"azurerm_eventhub":                              eventhub.RegisterDataSource(),
		"azurerm_eventhub_authorization_rule":           eventhub_authorization_rule.RegisterDataSource(),
		"azurerm_eventhub_cluster":                      eventhub_cluster.RegisterDataSource(),
		"azurerm_eventhub_consumer_group":               eventhub_consumer_group.RegisterDataSource(),
		"azurerm_eventhub_namespace":                    eventhub_namespace.RegisterDataSource(),
		"azurerm_eventhub_namespace_authorization_rule": eventhub_namespace_authorization_rule.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventhub":                                    eventhub.RegisterResource(),
		"azurerm_eventhub_authorization_rule":                 eventhub_authorization_rule.RegisterResource(),
		"azurerm_eventhub_cluster":                            eventhub_cluster.RegisterResource(),
		"azurerm_eventhub_namespace":                          eventhub_namespace.RegisterResource(),
		"azurerm_eventhub_namespace_authorization_rule":       eventhub_namespace_authorization_rule.RegisterResource(),
		"azurerm_eventhub_namespace_customer_managed_key":     eventhub_namespace_customer_managed_key.RegisterResource(),
		"azurerm_eventhub_namespace_disaster_recovery_config": eventhub_namespace_disaster_recovery_config.RegisterResource(),
		"azurerm_eventhub_namespace_schema_group":             eventhub_namespace_schema_group.RegisterResource(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		eventhub_sas.DataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		eventhub_consumer_group.Resource{},
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
		eventhub_consumer_group.ListResource{},
		eventhub.ListResource{},
	}
}
