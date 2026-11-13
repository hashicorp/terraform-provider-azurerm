// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/armdisasterrecoveries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/authorizationrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/consumergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/ehnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/networkrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2026-01-01/schemagroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ClusterClient                          *clusters.ClustersClient
	ConsumerGroupClient                    *consumergroups.ConsumerGroupsClient
	DisasterRecoveryConfigsClient          *armdisasterrecoveries.ArmDisasterRecoveriesClient
	DisasterRecoveryNameAvailabilityClient *ehnamespaces.EHNamespacesClient
	EventHubsClient                        *eventhubs.EventHubsClient
	EventHubAuthorizationRulesClient       *eventhubs.EventHubsClient
	NamespacesClient                       *ehnamespaces.EHNamespacesClient
	NamespaceAuthorizationRulesClient      *authorizationrules.AuthorizationRulesClient
	NetworkRuleSetsClient                  *networkrulesets.NetworkRuleSetsClient
	SchemaRegistryClient                   *schemagroups.SchemaGroupsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	clustersClient, err := clusters.NewClustersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Clusters Client: %+v", err)
	}
	o.Configure(clustersClient.Client, o.Authorizers.ResourceManager)

	consumerGroupsClient, err := consumergroups.NewConsumerGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ConsumerGroups Client: %+v", err)
	}
	o.Configure(consumerGroupsClient.Client, o.Authorizers.ResourceManager)

	disasterRecoveryConfigsClient, err := armdisasterrecoveries.NewArmDisasterRecoveriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DisasterRecoveryConfigs Client: %+v", err)
	}
	o.Configure(disasterRecoveryConfigsClient.Client, o.Authorizers.ResourceManager)

	disasterRecoveryNameAvailabilityClient, err := ehnamespaces.NewEHNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building DisasterRecoveryNameAvailability Client: %+v", err)
	}
	o.Configure(disasterRecoveryNameAvailabilityClient.Client, o.Authorizers.ResourceManager)

	eventhubsClient, err := eventhubs.NewEventHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building EventHubs Client: %+v", err)
	}
	o.Configure(eventhubsClient.Client, o.Authorizers.ResourceManager)

	eventHubAuthorizationRulesClient, err := eventhubs.NewEventHubsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building EventHubAuthorizationRules Client: %+v", err)
	}
	o.Configure(eventHubAuthorizationRulesClient.Client, o.Authorizers.ResourceManager)

	namespacesClient, err := ehnamespaces.NewEHNamespacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Namespaces Client: %+v", err)
	}
	o.Configure(namespacesClient.Client, o.Authorizers.ResourceManager)

	namespaceAuthorizationRulesClient, err := authorizationrules.NewAuthorizationRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NamespaceAuthorizationRules Client: %+v", err)
	}
	o.Configure(namespaceAuthorizationRulesClient.Client, o.Authorizers.ResourceManager)

	networkRuleSetsClient, err := networkrulesets.NewNetworkRuleSetsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building NetworkRuleSets Client: %+v", err)
	}
	o.Configure(networkRuleSetsClient.Client, o.Authorizers.ResourceManager)

	schemaRegistryClient, err := schemagroups.NewSchemaGroupsClientWithBaseURI(o.Environment.ResourceManager)
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
