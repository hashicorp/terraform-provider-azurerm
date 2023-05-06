package client

import (
	"github.com/Azure/go-autorest/autorest"
	loadtestservice_v2022_12_01 "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

func NewClient(o *common.ClientOptions) (*loadtestservice_v2022_12_01.Client, error) {
	client := loadtestservice_v2022_12_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		c.Authorizer = o.ResourceManagerAuthorizer
	})
	return &client, nil
}
