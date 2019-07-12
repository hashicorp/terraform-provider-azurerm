package eventgrid

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DomainsClient            eventgrid.DomainsClient
	EventSubscriptionsClient eventgrid.EventSubscriptionsClient
	TopicsClient             eventgrid.TopicsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.DomainsClient = eventgrid.NewDomainsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DomainsClient.Client, o.ResourceManagerAuthorizer)

	c.EventSubscriptionsClient = eventgrid.NewEventSubscriptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.EventSubscriptionsClient.Client, o.ResourceManagerAuthorizer)

	c.TopicsClient = eventgrid.NewTopicsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.TopicsClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
