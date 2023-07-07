package client

import (
	"github.com/Azure/go-autorest/autorest"
	azurestackhci_v2023_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *azurestackhci_v2023_03_01.Client {
	client := azurestackhci_v2023_03_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		o.ConfigureClient(c, o.ResourceManagerAuthorizer)
	})

	return &client
}
