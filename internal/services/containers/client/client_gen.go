package client

import (
	"github.com/Azure/go-autorest/autorest"
	containerservice_v2022_09_02_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*containerservice_v2022_09_02_preview.Client, error) {
	client := containerservice_v2022_09_02_preview.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client, nil
}
