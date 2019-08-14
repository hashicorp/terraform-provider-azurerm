package postgres

import (
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient      *postgresql.ConfigurationsClient
	DatabasesClient           *postgresql.DatabasesClient
	FirewallRulesClient       *postgresql.FirewallRulesClient
	ServersClient             *postgresql.ServersClient
	VirtualNetworkRulesClient *postgresql.VirtualNetworkRulesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ConfigurationsClient := postgresql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := postgresql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := postgresql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	ServersClient := postgresql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServersClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworkRulesClient := postgresql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:      &ConfigurationsClient,
		DatabasesClient:           &DatabasesClient,
		FirewallRulesClient:       &FirewallRulesClient,
		ServersClient:             &ServersClient,
		VirtualNetworkRulesClient: &VirtualNetworkRulesClient,
	}
}
