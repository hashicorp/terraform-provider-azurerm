package client

import (
	"github.com/Azure/go-autorest/autorest"
	managedidentity_v2022_01_31_preview "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2022-01-31-preview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *managedidentity_v2022_01_31_preview.Client {
	client := managedidentity_v2022_01_31_preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client
}
