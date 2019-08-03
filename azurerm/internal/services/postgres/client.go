package postgres

import (
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient      postgresql.ConfigurationsClient
	DatabasesClient           postgresql.DatabasesClient
	FirewallRulesClient       postgresql.FirewallRulesClient
	ServersClient             postgresql.ServersClient
	VirtualNetworkRulesClient postgresql.VirtualNetworkRulesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ConfigurationsClient = postgresql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	c.DatabasesClient = postgresql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatabasesClient.Client, o.ResourceManagerAuthorizer)

	c.FirewallRulesClient = postgresql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	c.ServersClient = postgresql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServersClient.Client, o.ResourceManagerAuthorizer)

	c.VirtualNetworkRulesClient = postgresql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
