package client

import (
	"github.com/Azure/go-autorest/autorest"
	media_v2021_11_01 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01"
	media_v2022_08_01 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	V20211101Client *media_v2021_11_01.Client
	V20220801Client *media_v2022_08_01.Client
}

func NewClient(o *common.ClientOptions) *Client {
	V20211101Client := media_v2021_11_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	V20220801Client := media_v2022_08_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &Client{
		V20211101Client: &V20211101Client,
		V20220801Client: &V20220801Client,
	}
}
