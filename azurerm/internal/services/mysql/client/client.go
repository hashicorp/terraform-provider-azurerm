package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient              *mysql.ConfigurationsClient
	DatabasesClient                   *mysql.DatabasesClient
	FirewallRulesClient               *mysql.FirewallRulesClient
	ServersClient                     *mysql.ServersClient
	ServerKeysClient                  *mysql.ServerKeysClient
	ServerSecurityAlertPoliciesClient *mysql.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient         *mysql.VirtualNetworkRulesClient
	ServerAdministratorsClient        *mysql.ServerAdministratorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ConfigurationsClient := mysql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := mysql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := mysql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	ServersClient := mysql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServersClient.Client, o.ResourceManagerAuthorizer)

	ServerKeysClient := mysql.NewServerKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerKeysClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := mysql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworkRulesClient := mysql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	serverAdministratorsClient := mysql.NewServerAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:              &ConfigurationsClient,
		DatabasesClient:                   &DatabasesClient,
		FirewallRulesClient:               &FirewallRulesClient,
		ServersClient:                     &ServersClient,
		ServerKeysClient:                  &ServerKeysClient,
		ServerSecurityAlertPoliciesClient: &serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:         &VirtualNetworkRulesClient,
		ServerAdministratorsClient:        &serverAdministratorsClient,
	}
}
