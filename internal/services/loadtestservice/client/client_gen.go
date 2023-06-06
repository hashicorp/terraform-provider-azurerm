package client

import (
	"github.com/Azure/go-autorest/autorest"
	loadtestserviceV20211201Preview "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20211201Preview loadtestserviceV20211201Preview.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20211201PreviewClient := loadtestserviceV20211201Preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &AutoClient{
		V20211201Preview: v20211201PreviewClient,
	}, nil
}
