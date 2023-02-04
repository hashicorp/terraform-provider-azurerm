package v2021_06_01

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2021-06-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2021-06-01/patchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2021-06-01/privateendpointconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2021-06-01/privatelinkresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2021-06-01/redis"
)

type Client struct {
	FirewallRules              *firewallrules.FirewallRulesClient
	PatchSchedules             *patchschedules.PatchSchedulesClient
	PrivateEndpointConnections *privateendpointconnections.PrivateEndpointConnectionsClient
	PrivateLinkResources       *privatelinkresources.PrivateLinkResourcesClient
	Redis                      *redis.RedisClient
}

func NewClientWithBaseURI(endpoint string, configureAuthFunc func(c *autorest.Client)) Client {

	firewallRulesClient := firewallrules.NewFirewallRulesClientWithBaseURI(endpoint)
	configureAuthFunc(&firewallRulesClient.Client)

	patchSchedulesClient := patchschedules.NewPatchSchedulesClientWithBaseURI(endpoint)
	configureAuthFunc(&patchSchedulesClient.Client)

	privateEndpointConnectionsClient := privateendpointconnections.NewPrivateEndpointConnectionsClientWithBaseURI(endpoint)
	configureAuthFunc(&privateEndpointConnectionsClient.Client)

	privateLinkResourcesClient := privatelinkresources.NewPrivateLinkResourcesClientWithBaseURI(endpoint)
	configureAuthFunc(&privateLinkResourcesClient.Client)

	redisClient := redis.NewRedisClientWithBaseURI(endpoint)
	configureAuthFunc(&redisClient.Client)

	return Client{
		FirewallRules:              &firewallRulesClient,
		PatchSchedules:             &patchSchedulesClient,
		PrivateEndpointConnections: &privateEndpointConnectionsClient,
		PrivateLinkResources:       &privateLinkResourcesClient,
		Redis:                      &redisClient,
	}
}
