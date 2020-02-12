package eventhub

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "EventHub"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventhub_consumer_group":               dataSourceEventHubConsumerGroup(),
		"azurerm_eventhub_namespace":                    dataSourceEventHubNamespace(),
		"azurerm_eventhub_namespace_authorization_rule": dataSourceEventHubNamespaceAuthorizationRule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventhub_authorization_rule":                 resourceArmEventHubAuthorizationRule(),
		"azurerm_eventhub_consumer_group":                     resourceArmEventHubConsumerGroup(),
		"azurerm_eventhub_namespace_authorization_rule":       resourceArmEventHubNamespaceAuthorizationRule(),
		"azurerm_eventhub_namespace_disaster_recovery_config": resourceArmEventHubNamespaceDisasterRecoveryConfig(),
		"azurerm_eventhub_namespace":                          resourceArmEventHubNamespace(),
		"azurerm_eventhub":                                    resourceArmEventHub()}
}
