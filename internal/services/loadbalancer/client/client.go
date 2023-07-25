// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type Client struct {
	LoadBalancersClient                   *network.LoadBalancersClient
	LoadBalancerBackendAddressPoolsClient *network.LoadBalancerBackendAddressPoolsClient
	LoadBalancingRulesClient              *network.LoadBalancerLoadBalancingRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	loadBalancersClient := network.NewLoadBalancersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loadBalancersClient.Client, o.ResourceManagerAuthorizer)

	loadBalancerBackendAddressPoolsClient := network.NewLoadBalancerBackendAddressPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loadBalancerBackendAddressPoolsClient.Client, o.ResourceManagerAuthorizer)

	loadBalancingRulesClient := network.NewLoadBalancerLoadBalancingRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&loadBalancingRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		LoadBalancersClient:                   &loadBalancersClient,
		LoadBalancerBackendAddressPoolsClient: &loadBalancerBackendAddressPoolsClient,
		LoadBalancingRulesClient:              &loadBalancingRulesClient,
	}
}
