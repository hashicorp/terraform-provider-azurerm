package client

import (
	"github.com/Azure/go-autorest/autorest"
	nginx "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2022-08-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *nginx.Client {
	client := nginx.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client
}
