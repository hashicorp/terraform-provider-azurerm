package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	PoolsClient *pools.PoolsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	poolsClient, err := pools.NewPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Pools client: %+v", err)
	}
	o.Configure(poolsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		PoolsClient: poolsClient,
	}, nil
}
