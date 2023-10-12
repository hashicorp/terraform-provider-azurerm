package client

import (
	"fmt"

	devcenterV20230401 "github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20230401 devcenterV20230401.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20230401Client, err := devcenterV20230401.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for devcenter V20230401: %+v", err)
	}

	return &AutoClient{
		V20230401: *v20230401Client,
	}, nil
}
