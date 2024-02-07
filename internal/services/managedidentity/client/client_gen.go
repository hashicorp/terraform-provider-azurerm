package client

import (
	"fmt"

	managedidentityV20230131 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2023-01-31"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20230131 managedidentityV20230131.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20230131Client, err := managedidentityV20230131.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for managedidentity V20230131: %+v", err)
	}

	return &AutoClient{
		V20230131: *v20230131Client,
	}, nil
}
