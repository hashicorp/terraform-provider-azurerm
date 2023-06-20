package client

import (
	"github.com/Azure/go-autorest/autorest"
	fluidrelay_2022_05_26 "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *fluidrelay_2022_05_26.Client {
	client := fluidrelay_2022_05_26.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client
}
