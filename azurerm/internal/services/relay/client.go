package relay

import "github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"

type Client struct {
	NamespacesClient relay.NamespacesClient
}
