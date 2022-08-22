package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queuesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topicsauthorizationrule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DisasterRecoveryConfigsClient *disasterrecoveryconfigs.DisasterRecoveryConfigsClient
	NamespacesAuthClient          *namespacesauthorizationrule.NamespacesAuthorizationRuleClient
	NamespacesClient              *namespaces.NamespacesClient
	QueuesAuthClient              *queuesauthorizationrule.QueuesAuthorizationRuleClient
	QueuesClient                  *queues.QueuesClient
	SubscriptionsClient           *subscriptions.SubscriptionsClient
	SubscriptionRulesClient       *rules.RulesClient
	TopicsAuthClient              *topicsauthorizationrule.TopicsAuthorizationRuleClient
	TopicsClient                  *topics.TopicsClient
}

func NewClient(o *common.ClientOptions) *Client {
	DisasterRecoveryConfigsClient := disasterrecoveryconfigs.NewDisasterRecoveryConfigsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DisasterRecoveryConfigsClient.Client, o.ResourceManagerAuthorizer)

	NamespacesAuthClient := namespacesauthorizationrule.NewNamespacesAuthorizationRuleClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&NamespacesAuthClient.Client, o.ResourceManagerAuthorizer)

	NamespacesClient := namespaces.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&NamespacesClient.Client, o.ResourceManagerAuthorizer)

	QueuesAuthClient := queuesauthorizationrule.NewQueuesAuthorizationRuleClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&QueuesAuthClient.Client, o.ResourceManagerAuthorizer)

	QueuesClient := queues.NewQueuesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&QueuesClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionsClient := subscriptions.NewSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&SubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	SubscriptionRulesClient := rules.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&SubscriptionRulesClient.Client, o.ResourceManagerAuthorizer)

	TopicsAuthClient := topicsauthorizationrule.NewTopicsAuthorizationRuleClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&TopicsAuthClient.Client, o.ResourceManagerAuthorizer)

	TopicsClient := topics.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&TopicsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DisasterRecoveryConfigsClient: &DisasterRecoveryConfigsClient,
		NamespacesAuthClient:          &NamespacesAuthClient,
		NamespacesClient:              &NamespacesClient,
		QueuesAuthClient:              &QueuesAuthClient,
		QueuesClient:                  &QueuesClient,
		SubscriptionsClient:           &SubscriptionsClient,
		SubscriptionRulesClient:       &SubscriptionRulesClient,
		TopicsAuthClient:              &TopicsAuthClient,
		TopicsClient:                  &TopicsClient,
	}
}
