package client

import (
	"github.com/Azure/go-autorest/autorest"
	redis_2022_06_01 "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2022-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) *redis_2022_06_01.Client {
	client := redis_2022_06_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client
}
