package client

import (
	"fmt"

	containerservice_2024_04_01 "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01"
	containerservice_2025_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-03-01"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20240401 containerservice_2024_04_01.Client
	V20250301 containerservice_2025_03_01.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {
	v20240401Client, err := containerservice_2024_04_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for containerservice V20240401: %+v", err)
	}

	v20250301Client, err := containerservice_2025_03_01.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for containerservice V20250301: %+v", err)
	}

	return &AutoClient{
		V20240401: *v20240401Client,
		V20250301: *v20250301Client,
	}, nil
}
