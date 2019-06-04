package eventhub

import (
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
)

type Client struct {
	ConsumerGroupClient eventhub.ConsumerGroupsClient
	EventHubsClient     eventhub.EventHubsClient
	NamespacesClient    eventhub.NamespacesClient
}
