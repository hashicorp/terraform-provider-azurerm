package servicebus

import "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"

type Client struct {
	QueuesClient            servicebus.QueuesClient
	NamespacesClient        servicebus.NamespacesClient
	TopicsClient            servicebus.TopicsClient
	SubscriptionsClient     servicebus.SubscriptionsClient
	SubscriptionRulesClient servicebus.RulesClient
}
