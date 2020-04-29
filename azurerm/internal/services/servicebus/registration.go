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
		"azurerm_servicebus_namespace":                    dataSourceArmServiceBusNamespace(),
		"azurerm_servicebus_namespace_authorization_rule": dataSourceArmServiceBusNamespaceAuthorizationRule(),
		"azurerm_servicebus_topic_authorization_rule":     dataSourceArmServiceBusTopicAuthorizationRule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_servicebus_namespace":                    resourceArmServiceBusNamespace(),
		"azurerm_servicebus_namespace_authorization_rule": resourceArmServiceBusNamespaceAuthorizationRule(),
		"azurerm_servicebus_namespace_network_rule_set":   resourceServiceBusNamespaceNetworkRuleSet(),
		"azurerm_servicebus_queue":                        resourceArmServiceBusQueue(),
		"azurerm_servicebus_queue_authorization_rule":     resourceArmServiceBusQueueAuthorizationRule(),
		"azurerm_servicebus_subscription":                 resourceArmServiceBusSubscription(),
		"azurerm_servicebus_subscription_rule":            resourceArmServiceBusSubscriptionRule(),
		"azurerm_servicebus_topic_authorization_rule":     resourceArmServiceBusTopicAuthorizationRule(),
		"azurerm_servicebus_topic":                        resourceArmServiceBusTopic(),
	}
}
