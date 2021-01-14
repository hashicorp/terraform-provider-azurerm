package eventhub

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventhub":                              dataSourceEventHub(),
		"azurerm_eventhub_authorization_rule":           EventHubAuthorizationRuleDataSource(),
		"azurerm_eventhub_consumer_group":               EventHubConsumerGroupDataSource(),
		"azurerm_eventhub_namespace":                    EventHubNamespaceDataSource(),
		"azurerm_eventhub_namespace_authorization_rule": EventHubNamespaceDataSourceAuthorizationRule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventhub_authorization_rule":                 resourceEventHubAuthorizationRule(),
		"azurerm_eventhub_cluster":                            resourceEventHubCluster(),
		"azurerm_eventhub_namespace_authorization_rule":       resourceEventHubNamespaceAuthorizationRule(),
		"azurerm_eventhub_namespace_disaster_recovery_config": resourceEventHubNamespaceDisasterRecoveryConfig(),
		"azurerm_eventhub_namespace":                          resourceEventHubNamespace(),
		"azurerm_eventhub":                                    resourceEventHub(),
	}
}

// PackagePath is the relative path to this package
func (r Registration) PackagePath() string {
	return "TODO"
}

// DataSources returns a list of Data Sources supported by this Service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns a list of Resources supported by this Service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ConsumerGroupResource{},
	}
}
