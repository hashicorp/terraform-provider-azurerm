// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventhub

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
		"azurerm_eventhub":                              dataSourceEventHub(),
		"azurerm_eventhub_cluster":                      dataSourceEventHubCluster(),
		"azurerm_eventhub_authorization_rule":           EventHubAuthorizationRuleDataSource(),
		"azurerm_eventhub_consumer_group":               EventHubConsumerGroupDataSource(),
		"azurerm_eventhub_namespace":                    EventHubNamespaceDataSource(),
		"azurerm_eventhub_namespace_authorization_rule": EventHubNamespaceDataSourceAuthorizationRule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_eventhub_authorization_rule":                 resourceEventHubAuthorizationRule(),
		"azurerm_eventhub_cluster":                            resourceEventHubCluster(),
		"azurerm_eventhub_namespace_authorization_rule":       resourceEventHubNamespaceAuthorizationRule(),
		"azurerm_eventhub_namespace_customer_managed_key":     resourceEventHubNamespaceCustomerManagedKey(),
		"azurerm_eventhub_namespace_disaster_recovery_config": resourceEventHubNamespaceDisasterRecoveryConfig(),
		"azurerm_eventhub_namespace":                          resourceEventHubNamespace(),
		"azurerm_eventhub_namespace_schema_group":             resourceEventHubNamespaceSchemaRegistry(),
		"azurerm_eventhub":                                    resourceEventHub(),
	}
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		EventHubSharedAccessSignatureDataSource{},
	}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ConsumerGroupResource{},
	}
}
