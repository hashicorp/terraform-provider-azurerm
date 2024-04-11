package client

import (
	"fmt"

	containerserviceV20230302Preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview"
	containerserviceV20231015 "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-10-15"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type AutoClient struct {
	V20230302Preview containerserviceV20230302Preview.Client
	V20231015        containerserviceV20231015.Client
}

func NewClient(o *common.ClientOptions) (*AutoClient, error) {

	v20230302PreviewClient, err := containerserviceV20230302Preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for containerservice V20230302Preview: %+v", err)
	}

	v20231015Client, err := containerserviceV20231015.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building client for containerservice V20231015: %+v", err)
	}

	return &AutoClient{
		V20230302Preview: *v20230302PreviewClient,
		V20231015:        *v20231015Client,
	}, nil
}
