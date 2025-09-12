// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/patchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Redis          *redis.RedisClient
	FirewallRules  *firewallrules.FirewallRulesClient
	PatchSchedules *patchschedules.PatchSchedulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	redisClient, err := redis.NewRedisClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building clients for Redis: %+v", err)
	}
	o.Configure(redisClient.Client, o.Authorizers.ResourceManager)

	fireWallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building clients for Redis: %+v", err)
	}
	o.Configure(fireWallRulesClient.Client, o.Authorizers.ResourceManager)

	patchSchedulesClient, err := patchschedules.NewPatchSchedulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building clients for Redis: %+v", err)
	}
	o.Configure(patchSchedulesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		Redis:          redisClient,
		FirewallRules:  fireWallRulesClient,
		PatchSchedules: patchSchedulesClient,
	}, nil
}
