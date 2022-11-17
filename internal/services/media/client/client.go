package client

import (
	"github.com/Azure/go-autorest/autorest"
	media_v2020_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01"
	media_v2021_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-05-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	V20200501Client *media_v2020_05_01.Client
	V20210501Client *media_v2021_05_01.Client
}

func NewClient(o *common.ClientOptions) *Client {
	v2020Client := media_v2020_05_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	v2021Client := media_v2021_05_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &Client{
		V20200501Client: &v2020Client,
		V20210501Client: &v2021Client,
	}
}
