package redis

import (
	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	Client               *redis.Client
	FirewallRulesClient  *redis.FirewallRulesClient
	PatchSchedulesClient *redis.PatchSchedulesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	client := redis.NewClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&client.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := redis.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	PatchSchedulesClient := redis.NewPatchSchedulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PatchSchedulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		Client:               &client,
		FirewallRulesClient:  &FirewallRulesClient,
		PatchSchedulesClient: &PatchSchedulesClient,
	}
}
