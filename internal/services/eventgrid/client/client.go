package client

import (
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DomainsClient                       *eventgrid.DomainsClient
	DomainTopicsClient                  *eventgrid.DomainTopicsClient
	EventSubscriptionsClient            *eventgrid.EventSubscriptionsClient
	TopicsClient                        *eventgrid.TopicsClient
	SystemTopicsClient                  *eventgrid.SystemTopicsClient
	SystemTopicEventSubscriptionsClient *eventgrid.SystemTopicEventSubscriptionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	DomainsClient := eventgrid.NewDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DomainsClient.Client, o.ResourceManagerAuthorizer)

	DomainTopicsClient := eventgrid.NewDomainTopicsClient(o.SubscriptionId)
	o.ConfigureClient(&DomainTopicsClient.Client, o.ResourceManagerAuthorizer)

	EventSubscriptionsClient := eventgrid.NewEventSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EventSubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	TopicsClient := eventgrid.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TopicsClient.Client, o.ResourceManagerAuthorizer)

	SystemTopicsClient := eventgrid.NewSystemTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SystemTopicsClient.Client, o.ResourceManagerAuthorizer)

	SystemTopicEventSubscriptionsClient := eventgrid.NewSystemTopicEventSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SystemTopicEventSubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DomainsClient:                       &DomainsClient,
		EventSubscriptionsClient:            &EventSubscriptionsClient,
		DomainTopicsClient:                  &DomainTopicsClient,
		TopicsClient:                        &TopicsClient,
		SystemTopicsClient:                  &SystemTopicsClient,
		SystemTopicEventSubscriptionsClient: &SystemTopicEventSubscriptionsClient,
	}
}
