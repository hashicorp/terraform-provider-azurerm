// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/namespacesauthorizationrule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/queues"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/subscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2026-01-01/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ArmDisasterRecoveriesClient   *armdisasterrecoveries.ArmDisasterRecoveriesClient
	DisasterRecoveryConfigsClient *disasterrecoveryconfigs.DisasterRecoveryConfigsClient
	NamespacesAuthClient          *namespacesauthorizationrule.NamespacesAuthorizationRuleClient
	NamespacesClient              *namespaces.NamespacesClient
	QueuesClient                  *queues.QueuesClient
	SubscriptionsClient           *subscriptions.SubscriptionsClient
	SubscriptionRulesClient       *rules.RulesClient
	TopicsClient                  *topics.TopicsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	armDisasterRecoveriesClient, err := armdisasterrecoveries.NewArmDisasterRecoveriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ArmDisasterRecoveries client: %+v", err)
	}
	o.Configure(armDisasterRecoveriesClient.Client, o.Authorizers.ResourceManager)

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

	topicsClient, err := topics.NewTopicsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Topics client: %+v", err)
	}
	o.Configure(topicsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ArmDisasterRecoveriesClient:   armDisasterRecoveriesClient,
		DisasterRecoveryConfigsClient: disasterRecoveryConfigsClient,
		NamespacesAuthClient:          namespacesAuthClient,
		NamespacesClient:              namespacesClient,
		QueuesClient:                  queuesClient,
		SubscriptionsClient:           subscriptionsClient,
		SubscriptionRulesClient:       subscriptionRulesClient,
		TopicsClient:                  topicsClient,
	}, nil
}
