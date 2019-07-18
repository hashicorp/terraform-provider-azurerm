package devspace

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ControllersClient devspaces.ControllersClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ControllersClient = devspaces.NewControllersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ControllersClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
