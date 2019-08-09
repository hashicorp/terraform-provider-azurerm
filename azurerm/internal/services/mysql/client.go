package mysql

import (
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient      mysql.ConfigurationsClient
	DatabasesClient           mysql.DatabasesClient
	FirewallRulesClient       mysql.FirewallRulesClient
	ServersClient             mysql.ServersClient
	VirtualNetworkRulesClient mysql.VirtualNetworkRulesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ConfigurationsClient = mysql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	c.DatabasesClient = mysql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatabasesClient.Client, o.ResourceManagerAuthorizer)

	c.FirewallRulesClient = mysql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	c.ServersClient = mysql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServersClient.Client, o.ResourceManagerAuthorizer)

	c.VirtualNetworkRulesClient = mysql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
