package client

import (
	"github.com/Azure/go-autorest/autorest"
	loadtestserviceV20221201 "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20221201 loadtestserviceV20221201.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20221201Client := loadtestserviceV20221201.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &AutoClient{
		V20221201: v20221201Client,
	}, nil
}
