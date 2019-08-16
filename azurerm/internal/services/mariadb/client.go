package mariadb

import (
	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient *mariadb.ConfigurationsClient
	DatabasesClient      *mariadb.DatabasesClient
	FirewallRulesClient  *mariadb.FirewallRulesClient
	ServersClient        *mariadb.ServersClient
}

func BuildClient(o *common.ClientOptions) *Client {

	configurationsClient := mariadb.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := mariadb.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := mariadb.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	ServersClient := mariadb.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServersClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient: &configurationsClient,
		DatabasesClient:      &DatabasesClient,
		FirewallRulesClient:  &FirewallRulesClient,
		ServersClient:        &ServersClient,
	}
}
