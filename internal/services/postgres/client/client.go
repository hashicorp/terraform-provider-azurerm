// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	flexibleserveradministrators "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratormicrosoftentras"
	backupsautomaticandondemand "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/backupautomaticandondemands"
	flexibleserverconfigurations "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	flexibleserverdatabases "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/databases"
	flexibleserverfirewallrules "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/firewallrules"
	flexibleservers "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	flexibleservervirtualendpoints "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupsClient                       *backupsautomaticandondemand.BackupAutomaticAndOnDemandsClient
	FlexibleServersClient               *flexibleservers.ServersClient
	FlexibleServersConfigurationsClient *flexibleserverconfigurations.ConfigurationsClient
	FlexibleServerFirewallRuleClient    *flexibleserverfirewallrules.FirewallRulesClient
	FlexibleServerDatabaseClient        *flexibleserverdatabases.DatabasesClient
	FlexibleServerAdministratorsClient  *flexibleserveradministrators.AdministratorMicrosoftEntrasClient
	VirtualEndpointClient               *flexibleservervirtualendpoints.VirtualEndpointsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupsClient, err := backupsautomaticandondemand.NewBackupAutomaticAndOnDemandsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Backups client: %+v", err)
	}
	o.Configure(backupsClient.Client, o.Authorizers.ResourceManager)

	flexibleServersClient, err := flexibleservers.NewServersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServers client: %+v", err)
	}
	o.Configure(flexibleServersClient.Client, o.Authorizers.ResourceManager)

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

	flexibleServerAdministratorsClient, err := flexibleserveradministrators.NewAdministratorMicrosoftEntrasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerAdministrators client: %+v", err)
	}
	o.Configure(flexibleServerAdministratorsClient.Client, o.Authorizers.ResourceManager)

	virtualEndpointClient, err := flexibleservervirtualendpoints.NewVirtualEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerVirtualEndpoint client: %+v", err)
	}
	o.Configure(virtualEndpointClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BackupsClient:                       backupsClient,
		FlexibleServersConfigurationsClient: flexibleServerConfigurationsClient,
		FlexibleServersClient:               flexibleServersClient,
		FlexibleServerFirewallRuleClient:    flexibleServerFirewallRuleClient,
		FlexibleServerDatabaseClient:        flexibleServerDatabaseClient,
		FlexibleServerAdministratorsClient:  flexibleServerAdministratorsClient,
		VirtualEndpointClient:               virtualEndpointClient,
	}, nil
}
