package client

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

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := postgresql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	databasesClient := postgresql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databasesClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient := postgresql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	serversClient := postgresql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serversClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := postgresql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:      &configurationsClient,
		DatabasesClient:           &databasesClient,
		FirewallRulesClient:       &firewallRulesClient,
		ServersClient:             &serversClient,
		VirtualNetworkRulesClient: &virtualNetworkRulesClient,
	}
}
