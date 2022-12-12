package client

import (
	"github.com/Azure/go-autorest/autorest"
	dns_v2018_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *dns_v2018_05_01.Client {
	client := dns_v2018_05_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client
}
