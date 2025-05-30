package client

import (
	"fmt"

	containerservice_2024_04_01 "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20240401 containerservice_2024_04_01.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20240401Client, err := containerservice_2024_04_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for containerservice V20231015: %+v", err)
	}

	return &AutoClient{
		V20240401: *v20240401Client,
	}, nil
}
