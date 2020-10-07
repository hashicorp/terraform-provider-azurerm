package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2020-04-01-preview/eventgrid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DomainsClient            *eventgrid.DomainsClient
	DomainTopicsClient       *eventgrid.DomainTopicsClient
	EventSubscriptionsClient *eventgrid.EventSubscriptionsClient
	TopicsClient             *eventgrid.TopicsClient
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

	return &Client{
		DomainsClient:            &DomainsClient,
		EventSubscriptionsClient: &EventSubscriptionsClient,
		DomainTopicsClient:       &DomainTopicsClient,
		TopicsClient:             &TopicsClient,
	}
}
