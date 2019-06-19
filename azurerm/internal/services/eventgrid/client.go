package eventgrid

import "github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"

type Client struct {
	DomainsClient            eventgrid.DomainsClient
	EventSubscriptionsClient eventgrid.EventSubscriptionsClient
	TopicsClient             eventgrid.TopicsClient
}
