package client

import (
	"fmt"

	managedidentityV20241130 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2024-11-30"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20241130 managedidentityV20241130.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20241130Client, err := managedidentityV20241130.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for managedidentity V20230131: %+v", err)
	}

	return &AutoClient{
		V20241130: *v20241130Client,
	}, nil
}
