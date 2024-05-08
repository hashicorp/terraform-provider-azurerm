// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	LoadBalancersClient *loadbalancers.LoadBalancersClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	loadBalancersClient, err := loadbalancers.NewLoadBalancersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building loadBalancers client: %+v", err)
	}
	o.Configure(loadBalancersClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		LoadBalancersClient: loadBalancersClient,
	}, nil
}
