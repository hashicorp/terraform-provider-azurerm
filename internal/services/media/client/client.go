package client

import (
	"github.com/Azure/go-autorest/autorest"
	mediaV20200501 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01"
	mediaV20211101 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01"
	mediaV20220801 "github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	V20200501Client *mediaV20200501.Client
	V20211101Client *mediaV20211101.Client
	V20220801Client *mediaV20220801.Client
}

func NewClient(o *common.ClientOptions) *Client {
	V20200501Client := mediaV20200501.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	V20211101Client := mediaV20211101.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	V20220801Client := mediaV20220801.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &Client{
		V20200501Client: &V20200501Client,
		V20211101Client: &V20211101Client,
		V20220801Client: &V20220801Client,
	}
}
