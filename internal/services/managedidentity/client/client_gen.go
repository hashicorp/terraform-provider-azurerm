package client

import (
	"github.com/Azure/go-autorest/autorest"
	managedidentityV20220131Preview "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2022-01-31-preview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20220131Preview managedidentityV20220131Preview.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20220131PreviewClient := managedidentityV20220131Preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &AutoClient{
		V20220131Preview: v20220131PreviewClient,
	}, nil
}
