package client

import (
	"github.com/Azure/go-autorest/autorest"
	containerserviceV20220902Preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview"
	containerserviceV20230402Preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-04-02-preview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20220902Preview containerserviceV20220902Preview.Client
	V20230402Preview containerserviceV20230402Preview.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20220902PreviewClient := containerserviceV20220902Preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	v20230402PreviewClient := containerserviceV20230402Preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})

	return &AutoClient{
		V20220902Preview: v20220902PreviewClient,
		V20230402Preview: v20230402PreviewClient,
	}, nil
}
