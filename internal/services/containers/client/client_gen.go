package client

import (
	"github.com/Azure/go-autorest/autorest"
	containerserviceV20220902Preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview"
	containerserviceV20230302Preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview"
	containerserviceV20231015 "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20220902Preview containerserviceV20220902Preview.Client
	V20230302Preview containerserviceV20230302Preview.Client
	V20231015        containerserviceV20231015.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20220902PreviewClient := containerserviceV20220902Preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		o.ConfigureClient(c, o.ResourceManagerAuthorizer)
	})

	v20230302PreviewClient := containerserviceV20230302Preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		o.ConfigureClient(c, o.ResourceManagerAuthorizer)
	})

	v20231015Client := containerserviceV20231015.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		o.ConfigureClient(c, o.ResourceManagerAuthorizer)
	})

	return &AutoClient{
		V20220902Preview: v20220902PreviewClient,
		V20230302Preview: v20230302PreviewClient,
		V20231015:        v20231015Client,
	}, nil
}
