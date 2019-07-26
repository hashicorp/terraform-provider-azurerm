package redis

import (
	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client               redis.Client
	FirewallRulesClient  redis.FirewallRulesClient
	PatchSchedulesClient redis.PatchSchedulesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.Client = redis.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.Client.Client, o.ResourceManagerAuthorizer)

	c.FirewallRulesClient = redis.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	c.PatchSchedulesClient = redis.NewPatchSchedulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.PatchSchedulesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
