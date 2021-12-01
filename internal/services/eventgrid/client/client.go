package client

import (
	eventgridAlias "github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2021-12-01/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-10-15-preview/eventgrid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DomainsClient                       *eventgridAlias.DomainsClient
	DomainTopicsClient                  *eventgrid.DomainTopicsClient
	EventSubscriptionsClient            *eventgrid.EventSubscriptionsClient
	TopicsClient                        *eventgridAlias.TopicsClient
	SystemTopicsClient                  *eventgridAlias.SystemTopicsClient
	SystemTopicEventSubscriptionsClient *eventgrid.SystemTopicEventSubscriptionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	DomainsClient := eventgridAlias.NewDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DomainsClient.Client, o.ResourceManagerAuthorizer)

	DomainTopicsClient := eventgrid.NewDomainTopicsClient(o.SubscriptionId)
	o.ConfigureClient(&DomainTopicsClient.Client, o.ResourceManagerAuthorizer)

	EventSubscriptionsClient := eventgrid.NewEventSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&EventSubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	TopicsClient := eventgridAlias.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&TopicsClient.Client, o.ResourceManagerAuthorizer)

	SystemTopicsClient := eventgridAlias.NewSystemTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
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
