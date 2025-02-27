package client

import (
	"fmt"

	chaosstudioV20231101 "github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20231101 chaosstudioV20231101.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20231101Client, err := chaosstudioV20231101.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for chaosstudio V20231101: %+v", err)
	}

	return &AutoClient{
		V20231101: *v20231101Client,
	}, nil
}
