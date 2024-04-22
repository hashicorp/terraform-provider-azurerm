package client

import (
	"fmt"

	extendedLocation20210815 "github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	*extendedLocation20210815.Client
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	client, err := extendedLocation20210815.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, fmt.Errorf("building clients for Network: %+v", err)
	}

	return &Client{
		Client: client,
	}, nil
}
