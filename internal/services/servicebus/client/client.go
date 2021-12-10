package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/servicebus/mgmt/2021-06-01-preview/servicebus"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	QueuesClient                  *servicebus.QueuesClient
	DisasterRecoveryConfigsClient *servicebus.DisasterRecoveryConfigsClient
	NamespacesClient              *servicebus.NamespacesClient
	TopicsClient                  *servicebus.TopicsClient
	SubscriptionsClient           *servicebus.SubscriptionsClient
	SubscriptionRulesClient       *servicebus.RulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	QueuesClient := servicebus.NewQueuesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&QueuesClient.Client, o.ResourceManagerAuthorizer)

	DisasterRecoveryConfigsClient := servicebus.NewDisasterRecoveryConfigsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DisasterRecoveryConfigsClient.Client, o.ResourceManagerAuthorizer)

	NamespacesClient := servicebus.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	TopicsClient := servicebus.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TopicsClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionsClient := servicebus.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionRulesClient := servicebus.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SubscriptionRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		QueuesClient:                  &QueuesClient,
		DisasterRecoveryConfigsClient: &DisasterRecoveryConfigsClient,
		NamespacesClient:              &NamespacesClient,
		TopicsClient:                  &TopicsClient,
		SubscriptionsClient:           &SubscriptionsClient,
		SubscriptionRulesClient:       &SubscriptionRulesClient,
	}
}
