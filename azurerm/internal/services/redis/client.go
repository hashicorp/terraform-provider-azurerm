package redis

import "github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"

type Client struct {
	Client               redis.Client
	FirewallRulesClient  redis.FirewallRulesClient
	PatchSchedulesClient redis.PatchSchedulesClient
}
