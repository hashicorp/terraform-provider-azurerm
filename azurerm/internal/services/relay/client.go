package relay

import (
	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	NamespacesClient relay.NamespacesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.NamespacesClient = relay.NewNamespacesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.NamespacesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
