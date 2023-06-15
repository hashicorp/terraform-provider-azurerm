package client

import (
	"github.com/Azure/go-autorest/autorest"
	managedidentityV20230131 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2023-01-31"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20230131 managedidentityV20230131.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20230131Client := managedidentityV20230131.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &AutoClient{
		V20230131: v20230131Client,
	}, nil
}
