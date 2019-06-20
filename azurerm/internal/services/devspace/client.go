package devspace

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/devspaces/mgmt/2018-06-01-preview/devspaces"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ControllersClient devspaces.ControllersClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.ControllersClient = devspaces.NewControllersClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ControllersClient.Client, authorizer)

	return &c
}
