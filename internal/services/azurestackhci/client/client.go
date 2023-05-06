package client

import (
	"github.com/Azure/go-autorest/autorest"
	azurestackhci_v2022_12_01 "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2022-12-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *azurestackhci_v2022_12_01.Client {
	client := azurestackhci_v2022_12_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &client
}
