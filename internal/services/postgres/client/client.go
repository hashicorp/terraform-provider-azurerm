package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/replicas"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serveradministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/serversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/virtualnetworkrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2020-01-01/serverkeys"
	flexibleserverconfigurations "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/configurations"
	flexibleserverfirewallrules "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/serverrestart"
	flexibleserveradministrators "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/administrators"
	flexibleserverdatabases "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/databases"
	flexibleservers "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient                *configurations.ConfigurationsClient
	DatabasesClient                     *databases.DatabasesClient
	FirewallRulesClient                 *firewallrules.FirewallRulesClient
	FlexibleServersClient               *flexibleservers.ServersClient
	FlexibleServersConfigurationsClient *flexibleserverconfigurations.ConfigurationsClient
	FlexibleServerFirewallRuleClient    *flexibleserverfirewallrules.FirewallRulesClient
	FlexibleServerDatabaseClient        *flexibleserverdatabases.DatabasesClient
	FlexibleServerAdministratorsClient  *flexibleserveradministrators.AdministratorsClient
	ServersClient                       *servers.ServersClient
	ServerRestartClient                 *serverrestart.ServerRestartClient
	ServerKeysClient                    *serverkeys.ServerKeysClient
	ServerSecurityAlertPoliciesClient   *serversecurityalertpolicies.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient           *virtualnetworkrules.VirtualNetworkRulesClient
	ServerAdministratorsClient          *serveradministrators.ServerAdministratorsClient
	ReplicasClient                      *replicas.ReplicasClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := configurations.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	databasesClient := databases.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&databasesClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient := firewallrules.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	serversClient := servers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serversClient.Client, o.ResourceManagerAuthorizer)

	serverKeysClient := serverkeys.NewServerKeysClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serverKeysClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := serversecurityalertpolicies.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	serverAdministratorsClient := serveradministrators.NewServerAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&serverAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	replicasClient := replicas.NewReplicasClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&replicasClient.Client, o.ResourceManagerAuthorizer)

	flexibleServersClient := flexibleservers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServersClient.Client, o.ResourceManagerAuthorizer)

	restartServerClient := serverrestart.NewServerRestartClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&restartServerClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerFirewallRuleClient := flexibleserverfirewallrules.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerFirewallRuleClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerDatabaseClient := flexibleserverdatabases.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerDatabaseClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerConfigurationsClient := flexibleserverconfigurations.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerAdministratorsClient := flexibleserveradministrators.NewAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:                &configurationsClient,
		DatabasesClient:                     &databasesClient,
		FirewallRulesClient:                 &firewallRulesClient,
		FlexibleServersConfigurationsClient: &flexibleServerConfigurationsClient,
		FlexibleServersClient:               &flexibleServersClient,
		ServerRestartClient:                 &restartServerClient,
		FlexibleServerFirewallRuleClient:    &flexibleServerFirewallRuleClient,
		FlexibleServerDatabaseClient:        &flexibleServerDatabaseClient,
		FlexibleServerAdministratorsClient:  &flexibleServerAdministratorsClient,
		ServersClient:                       &serversClient,
		ServerKeysClient:                    &serverKeysClient,
		ServerSecurityAlertPoliciesClient:   &serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:           &virtualNetworkRulesClient,
		ServerAdministratorsClient:          &serverAdministratorsClient,
		ReplicasClient:                      &replicasClient,
	}
}
