package mariadb

import (
	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient mariadb.ConfigurationsClient
	DatabasesClient      mariadb.DatabasesClient
	FirewallRulesClient  mariadb.FirewallRulesClient
	ServersClient        mariadb.ServersClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ConfigurationsClient = mariadb.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	c.DatabasesClient = mariadb.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatabasesClient.Client, o.ResourceManagerAuthorizer)

	c.FirewallRulesClient = mariadb.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	c.ServersClient = mariadb.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServersClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
