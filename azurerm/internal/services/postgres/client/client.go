package client

import (
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2020-01-01/postgresql"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2021-06-01/postgresqlflexibleservers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient                *postgresql.ConfigurationsClient
	DatabasesClient                     *postgresql.DatabasesClient
	FirewallRulesClient                 *postgresql.FirewallRulesClient
	FlexibleServersClient               *postgresqlflexibleservers.ServersClient
	FlexibleServersConfigurationsClient *postgresqlflexibleservers.ConfigurationsClient
	FlexibleServerFirewallRuleClient    *postgresqlflexibleservers.FirewallRulesClient
	ServersClient                       *postgresql.ServersClient
	ServerKeysClient                    *postgresql.ServerKeysClient
	ServerSecurityAlertPoliciesClient   *postgresql.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient           *postgresql.VirtualNetworkRulesClient
	ServerAdministratorsClient          *postgresql.ServerAdministratorsClient
	ReplicasClient                      *postgresql.ReplicasClient
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

	serverKeysClient := postgresql.NewServerKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverKeysClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := postgresql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := postgresql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	serverAdministratorsClient := postgresql.NewServerAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	replicasClient := postgresql.NewReplicasClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&replicasClient.Client, o.ResourceManagerAuthorizer)

	flexibleServersClient := postgresqlflexibleservers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&flexibleServersClient.Client, o.ResourceManagerAuthorizer)

	FlexibleServerFirewallRuleClient := postgresqlflexibleservers.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FlexibleServerFirewallRuleClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerConfigurationsClient := postgresqlflexibleservers.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&flexibleServerConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:                &configurationsClient,
		DatabasesClient:                     &databasesClient,
		FirewallRulesClient:                 &firewallRulesClient,
		FlexibleServersConfigurationsClient: &flexibleServerConfigurationsClient,
		FlexibleServersClient:               &flexibleServersClient,
		FlexibleServerFirewallRuleClient:    &FlexibleServerFirewallRuleClient,
		ServersClient:                       &serversClient,
		ServerKeysClient:                    &serverKeysClient,
		ServerSecurityAlertPoliciesClient:   &serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:           &virtualNetworkRulesClient,
		ServerAdministratorsClient:          &serverAdministratorsClient,
		ReplicasClient:                      &replicasClient,
	}
}
