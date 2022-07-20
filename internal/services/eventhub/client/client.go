package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/authorizationruleseventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/checknameavailabilitydisasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/consumergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/disasterrecoveryconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2017-04-01/eventhubs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2018-01-01-preview/eventhubsclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2018-01-01-preview/networkrulesets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-01-01-preview/namespaces"
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
}

func NewClient(o *common.ClientOptions) *Client {
	clustersClient := eventhubsclusters.NewEventHubsClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&clustersClient.Client, o.ResourceManagerAuthorizer)

	consumerGroupsClient := consumergroups.NewConsumerGroupsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&consumerGroupsClient.Client, o.ResourceManagerAuthorizer)

	disasterRecoveryConfigsClient := disasterrecoveryconfigs.NewDisasterRecoveryConfigsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&disasterRecoveryConfigsClient.Client, o.ResourceManagerAuthorizer)

	disasterRecoveryNameAvailabilityClient := checknameavailabilitydisasterrecoveryconfigs.NewCheckNameAvailabilityDisasterRecoveryConfigsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&disasterRecoveryNameAvailabilityClient.Client, o.ResourceManagerAuthorizer)

	eventhubsClient := eventhubs.NewEventHubsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&eventhubsClient.Client, o.ResourceManagerAuthorizer)

	eventHubAuthorizationRulesClient := authorizationruleseventhubs.NewAuthorizationRulesEventHubsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&eventHubAuthorizationRulesClient.Client, o.ResourceManagerAuthorizer)

	namespacesClient := namespaces.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&namespacesClient.Client, o.ResourceManagerAuthorizer)

	namespaceAuthorizationRulesClient := authorizationrulesnamespaces.NewAuthorizationRulesNamespacesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&namespaceAuthorizationRulesClient.Client, o.ResourceManagerAuthorizer)

	networkRuleSetsClient := networkrulesets.NewNetworkRuleSetsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&networkRuleSetsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ClusterClient:                          &clustersClient,
		ConsumerGroupClient:                    &consumerGroupsClient,
		DisasterRecoveryConfigsClient:          &disasterRecoveryConfigsClient,
		DisasterRecoveryNameAvailabilityClient: &disasterRecoveryNameAvailabilityClient,
		EventHubsClient:                        &eventhubsClient,
		EventHubAuthorizationRulesClient:       &eventHubAuthorizationRulesClient,
		NamespacesClient:                       &namespacesClient,
		NamespaceAuthorizationRulesClient:      &namespaceAuthorizationRulesClient,
		NetworkRuleSetsClient:                  &networkRuleSetsClient,
	}
}
