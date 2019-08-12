package subscription

import (
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-06-01/subscriptions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client subscriptions.Client
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.Client = subscriptions.NewClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&c.Client.Client, o.ResourceManagerAuthorizer)

	return &c
}
