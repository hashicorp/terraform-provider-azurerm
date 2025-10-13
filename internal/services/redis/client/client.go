// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/linkedserver"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/rediscacheaccesspolicyassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisfirewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	Redis                              *redis.RedisClient
	CacheAccessPoliciesClient          *rediscacheaccesspolicies.RedisCacheAccessPoliciesClient
	CacheAccessPolicyAssignmentsClient *rediscacheaccesspolicyassignments.RedisCacheAccessPolicyAssignmentsClient
	FirewallRulesClient                *redisfirewallrules.RedisFirewallRulesClient
	LinkedServerClient                 *linkedserver.LinkedServerClient
	PatchSchedulesClient               *redispatchschedules.RedisPatchSchedulesClient
	RedisResourcesClient               *redisresources.RedisResourcesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	redisClient, err := redis.NewRedisClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Redis Client: %+v", err)
	}
	o.Configure(redisClient.Client, o.Authorizers.ResourceManager)

	cacheAccessPoliciesClient, err := rediscacheaccesspolicies.NewRedisCacheAccessPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Cache Access Policies Client: %+v", err)
	}
	o.Configure(cacheAccessPoliciesClient.Client, o.Authorizers.ResourceManager)

	cacheAccessPolicyAssignmentsClient, err := rediscacheaccesspolicyassignments.NewRedisCacheAccessPolicyAssignmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Cache Access Policy Assignments Client: %+v", err)
	}
	o.Configure(cacheAccessPolicyAssignmentsClient.Client, o.Authorizers.ResourceManager)

	firewallRulesClient, err := redisfirewallrules.NewRedisFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall Rules Client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	linkedServerClient, err := linkedserver.NewLinkedServerClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Linked Server Client: %+v", err)
	}
	o.Configure(linkedServerClient.Client, o.Authorizers.ResourceManager)

	patchSchedulesClient, err := redispatchschedules.NewRedisPatchSchedulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Patch Schedules Client: %+v", err)
	}
	o.Configure(patchSchedulesClient.Client, o.Authorizers.ResourceManager)

	redisResourcesClient, err := redisresources.NewRedisResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Redis Resources Client: %+v", err)
	}
	o.Configure(redisResourcesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		Redis:                              redisClient,
		CacheAccessPoliciesClient:          cacheAccessPoliciesClient,
		CacheAccessPolicyAssignmentsClient: cacheAccessPolicyAssignmentsClient,
		FirewallRulesClient:                firewallRulesClient,
		LinkedServerClient:                 linkedServerClient,
		PatchSchedulesClient:               patchSchedulesClient,
		RedisResourcesClient:               redisResourcesClient,
	}, nil
}
