// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/namespacesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/queuesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/subscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2021-06-01-preview/topicsauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-01-01-preview/namespaces"
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

func NewClient(o *common.ClientOptions) (*Client, error) {
	disasterRecoveryConfigsClient, err := disasterrecoveryconfigs.NewDisasterRecoveryConfigsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DisasterRecoveryConfigs client: %+v", err)
	}
	o.Configure(disasterRecoveryConfigsClient.Client, o.Authorizers.ResourceManager)

	namespacesAuthClient, err := namespacesauthorizationrule.NewNamespacesAuthorizationRuleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NamespacesAuthorizationRule client: %+v", err)
	}
	o.Configure(namespacesAuthClient.Client, o.Authorizers.ResourceManager)

	namespacesClient, err := namespaces.NewNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Namespaces client: %+v", err)
	}
	o.Configure(namespacesClient.Client, o.Authorizers.ResourceManager)

	queuesAuthClient, err := queuesauthorizationrule.NewQueuesAuthorizationRuleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building QueuesAuthorizationRule client: %+v", err)
	}
	o.Configure(queuesAuthClient.Client, o.Authorizers.ResourceManager)

	queuesClient, err := queues.NewQueuesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Queues client: %+v", err)
	}
	o.Configure(queuesClient.Client, o.Authorizers.ResourceManager)

	subscriptionsClient, err := subscriptions.NewSubscriptionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Subscriptions client: %+v", err)
	}
	o.Configure(subscriptionsClient.Client, o.Authorizers.ResourceManager)

	subscriptionRulesClient, err := rules.NewRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Rules client: %+v", err)
	}
	o.Configure(subscriptionRulesClient.Client, o.Authorizers.ResourceManager)

	topicsAuthClient, err := topicsauthorizationrule.NewTopicsAuthorizationRuleClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building TopicsAuthorizationRule client: %+v", err)
	}
	o.Configure(topicsAuthClient.Client, o.Authorizers.ResourceManager)

	topicsClient, err := topics.NewTopicsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Topics client: %+v", err)
	}
	o.Configure(topicsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		DisasterRecoveryConfigsClient: disasterRecoveryConfigsClient,
		NamespacesAuthClient:          namespacesAuthClient,
		NamespacesClient:              namespacesClient,
		QueuesAuthClient:              queuesAuthClient,
		QueuesClient:                  queuesClient,
		SubscriptionsClient:           subscriptionsClient,
		SubscriptionRulesClient:       subscriptionRulesClient,
		TopicsAuthClient:              topicsAuthClient,
		TopicsClient:                  topicsClient,
	}, nil
}
