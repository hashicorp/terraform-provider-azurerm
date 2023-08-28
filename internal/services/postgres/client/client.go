// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2021-06-01/serverrestart"
	flexibleserveradministrators "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/administrators"
	flexibleserverdatabases "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/databases"
	flexibleserverfirewallrules "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/firewallrules"
	flexibleservers "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-03-01-preview/servers"
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

func NewClient(o *common.ClientOptions) (*Client, error) {
	configurationsClient, err := configurations.NewConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Configurations client: %+v", err)
	}
	o.Configure(configurationsClient.Client, o.Authorizers.ResourceManager)

	databasesClient, err := databases.NewDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Databases client: %+v", err)
	}
	o.Configure(databasesClient.Client, o.Authorizers.ResourceManager)

	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FirewallRules client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	serversClient, err := servers.NewServersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Servers client: %+v", err)
	}
	o.Configure(serversClient.Client, o.Authorizers.ResourceManager)

	serverKeysClient, err := serverkeys.NewServerKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServerKeys client: %+v", err)
	}
	o.Configure(serverKeysClient.Client, o.Authorizers.ResourceManager)

	serverSecurityAlertPoliciesClient, err := serversecurityalertpolicies.NewServerSecurityAlertPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServerSecurityAlertPolicies client: %+v", err)
	}
	o.Configure(serverSecurityAlertPoliciesClient.Client, o.Authorizers.ResourceManager)

	virtualNetworkRulesClient, err := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworkRules client: %+v", err)
	}
	o.Configure(virtualNetworkRulesClient.Client, o.Authorizers.ResourceManager)

	serverAdministratorsClient, err := serveradministrators.NewServerAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServerAdministrators client: %+v", err)
	}
	o.Configure(serverAdministratorsClient.Client, o.Authorizers.ResourceManager)

	replicasClient, err := replicas.NewReplicasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Replicas client: %+v", err)
	}
	o.Configure(replicasClient.Client, o.Authorizers.ResourceManager)

	flexibleServersClient, err := flexibleservers.NewServersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServers client: %+v", err)
	}
	o.Configure(flexibleServersClient.Client, o.Authorizers.ResourceManager)

	restartServerClient, err := serverrestart.NewServerRestartClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ServerRestart client: %+v", err)
	}
	o.Configure(restartServerClient.Client, o.Authorizers.ResourceManager)

	flexibleServerFirewallRuleClient, err := flexibleserverfirewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerFirewallRules client: %+v", err)
	}
	o.Configure(flexibleServerFirewallRuleClient.Client, o.Authorizers.ResourceManager)

	flexibleServerDatabaseClient, err := flexibleserverdatabases.NewDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerDatabases client: %+v", err)
	}
	o.Configure(flexibleServerDatabaseClient.Client, o.Authorizers.ResourceManager)

	flexibleServerConfigurationsClient, err := flexibleserverconfigurations.NewConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerConfigurations client: %+v", err)
	}
	o.Configure(flexibleServerConfigurationsClient.Client, o.Authorizers.ResourceManager)

	flexibleServerAdministratorsClient, err := flexibleserveradministrators.NewAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerAdministrators client: %+v", err)
	}
	o.Configure(flexibleServerAdministratorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ConfigurationsClient:                configurationsClient,
		DatabasesClient:                     databasesClient,
		FirewallRulesClient:                 firewallRulesClient,
		FlexibleServersConfigurationsClient: flexibleServerConfigurationsClient,
		FlexibleServersClient:               flexibleServersClient,
		ServerRestartClient:                 restartServerClient,
		FlexibleServerFirewallRuleClient:    flexibleServerFirewallRuleClient,
		FlexibleServerDatabaseClient:        flexibleServerDatabaseClient,
		FlexibleServerAdministratorsClient:  flexibleServerAdministratorsClient,
		ServersClient:                       serversClient,
		ServerKeysClient:                    serverKeysClient,
		ServerSecurityAlertPoliciesClient:   serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:           virtualNetworkRulesClient,
		ServerAdministratorsClient:          serverAdministratorsClient,
		ReplicasClient:                      replicasClient,
	}, nil
}
