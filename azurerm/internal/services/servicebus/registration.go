package servicebus

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_servicebus_namespace":                          dataSourceServiceBusNamespace(),
		"azurerm_servicebus_namespace_disaster_recovery_config": dataSourceServiceBusNamespaceDisasterRecoveryConfig(),
		"azurerm_servicebus_namespace_authorization_rule":       dataSourceServiceBusNamespaceAuthorizationRule(),
		"azurerm_servicebus_topic_authorization_rule":           dataSourceServiceBusTopicAuthorizationRule(),
		"azurerm_servicebus_queue_authorization_rule":           dataSourceServiceBusQueueAuthorizationRule(),
		"azurerm_servicebus_subscription":                       dataSourceServiceBusSubscription(),
		"azurerm_servicebus_topic":                              dataSourceServiceBusTopic(),
		"azurerm_servicebus_queue":                              dataSourceServiceBusQueue(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_servicebus_namespace":                          resourceServiceBusNamespace(),
		"azurerm_servicebus_namespace_disaster_recovery_config": resourceServiceBusNamespaceDisasterRecoveryConfig(),
		"azurerm_servicebus_namespace_authorization_rule":       resourceServiceBusNamespaceAuthorizationRule(),
		"azurerm_servicebus_namespace_network_rule_set":         resourceServiceBusNamespaceNetworkRuleSet(),
		"azurerm_servicebus_queue":                              resourceServiceBusQueue(),
		"azurerm_servicebus_queue_authorization_rule":           resourceServiceBusQueueAuthorizationRule(),
		"azurerm_servicebus_subscription":                       resourceServiceBusSubscription(),
		"azurerm_servicebus_subscription_rule":                  resourceServiceBusSubscriptionRule(),
		"azurerm_servicebus_topic_authorization_rule":           resourceServiceBusTopicAuthorizationRule(),
		"azurerm_servicebus_topic":                              resourceServiceBusTopic(),
	}
}
