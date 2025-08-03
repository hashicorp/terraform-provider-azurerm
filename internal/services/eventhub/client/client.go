// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationruleseventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/checknameavailabilitydisasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/consumergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/eventhubsclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/namespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/networkrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2024-01-01/schemaregistry"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClusterClient                          *eventhubsclusters.EventHubsClustersClient
	ConsumerGroupClient                    *consumergroups.ConsumerGroupsClient
	DisasterRecoveryConfigsClient          *disasterrecoveryconfigs.DisasterRecoveryConfigsClient
	DisasterRecoveryNameAvailabilityClient *checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityDisasterRecoveryConfigsClient
	EventHubsClient                        *eventhubs.EventHubsClient
	EventHubAuthorizationRulesClient       *authorizationruleseventhubs.AuthorizationRulesEventHubsClient
	NamespacesClient                       *namespaces.NamespacesClient
	NamespaceAuthorizationRulesClient      *authorizationrulesnamespaces.AuthorizationRulesNamespacesClient
	NetworkRuleSetsClient                  *networkrulesets.NetworkRuleSetsClient
	SchemaRegistryClient                   *schemaregistry.SchemaRegistryClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	clustersClient, err := eventhubsclusters.NewEventHubsClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters Client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	consumerGroupsClient, err := consumergroups.NewConsumerGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConsumerGroups Client: %+v", err)
	}
	o.Configure(consumerGroupsClient.Client, o.Authorizers.ResourceManager)

	disasterRecoveryConfigsClient, err := disasterrecoveryconfigs.NewDisasterRecoveryConfigsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DisasterRecoveryConfigs Client: %+v", err)
	}
	o.Configure(disasterRecoveryConfigsClient.Client, o.Authorizers.ResourceManager)

	disasterRecoveryNameAvailabilityClient, err := checknameavailabilitydisasterrecoveryconfigs.NewCheckNameAvailabilityDisasterRecoveryConfigsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DisasterRecoveryNameAvailability Client: %+v", err)
	}
	o.Configure(disasterRecoveryNameAvailabilityClient.Client, o.Authorizers.ResourceManager)

	eventhubsClient, err := eventhubs.NewEventHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building EventHubs Client: %+v", err)
	}
	o.Configure(eventhubsClient.Client, o.Authorizers.ResourceManager)

	eventHubAuthorizationRulesClient, err := authorizationruleseventhubs.NewAuthorizationRulesEventHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building EventHubAuthorizationRules Client: %+v", err)
	}
	o.Configure(eventHubAuthorizationRulesClient.Client, o.Authorizers.ResourceManager)

	namespacesClient, err := namespaces.NewNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Namspaces Client: %+v", err)
	}
	o.Configure(namespacesClient.Client, o.Authorizers.ResourceManager)

	namespaceAuthorizationRulesClient, err := authorizationrulesnamespaces.NewAuthorizationRulesNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NamespaceAuthorizationRules Client: %+v", err)
	}
	o.Configure(namespaceAuthorizationRulesClient.Client, o.Authorizers.ResourceManager)

	networkRuleSetsClient, err := networkrulesets.NewNetworkRuleSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NetworkRuleSets Client: %+v", err)
	}
	o.Configure(networkRuleSetsClient.Client, o.Authorizers.ResourceManager)

	schemaRegistryClient, err := schemaregistry.NewSchemaRegistryClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SchemaRegistry Client: %+v", err)
	}
	o.Configure(schemaRegistryClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ClusterClient:                          clustersClient,
		ConsumerGroupClient:                    consumerGroupsClient,
		DisasterRecoveryConfigsClient:          disasterRecoveryConfigsClient,
		DisasterRecoveryNameAvailabilityClient: disasterRecoveryNameAvailabilityClient,
		EventHubsClient:                        eventhubsClient,
		EventHubAuthorizationRulesClient:       eventHubAuthorizationRulesClient,
		NamespacesClient:                       namespacesClient,
		NamespaceAuthorizationRulesClient:      namespaceAuthorizationRulesClient,
		NetworkRuleSetsClient:                  networkRuleSetsClient,
		SchemaRegistryClient:                   schemaRegistryClient,
	}, nil
}
