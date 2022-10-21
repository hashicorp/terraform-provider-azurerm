package client

import (
	"github.com/Azure/go-autorest/autorest"
	managedidentity_v2018_11_30 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *managedidentity_v2018_11_30.Client {
	client := managedidentity_v2018_11_30.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client
}
