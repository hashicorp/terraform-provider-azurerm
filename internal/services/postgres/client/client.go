// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/administratormicrosoftentras"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/backupautomaticandondemands"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/virtualendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupsClient                       *backupautomaticandondemands.BackupAutomaticAndOnDemandsClient
	FlexibleServersClient               *servers.ServersClient
	FlexibleServersConfigurationsClient *configurations.ConfigurationsClient
	FlexibleServerFirewallRuleClient    *firewallrules.FirewallRulesClient
	FlexibleServerDatabaseClient        *databases.DatabasesClient
	FlexibleServerAdministratorsClient  *administratormicrosoftentras.AdministratorMicrosoftEntrasClient
	VirtualEndpointClient               *virtualendpoints.VirtualEndpointsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupsClient, err := backupautomaticandondemands.NewBackupAutomaticAndOnDemandsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Backups client: %+v", err)
	}
	o.Configure(backupsClient.Client, o.Authorizers.ResourceManager)

	flexibleServersClient, err := servers.NewServersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServers client: %+v", err)
	}
	o.Configure(flexibleServersClient.Client, o.Authorizers.ResourceManager)

	flexibleServerFirewallRuleClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerFirewallRules client: %+v", err)
	}
	o.Configure(flexibleServerFirewallRuleClient.Client, o.Authorizers.ResourceManager)

	flexibleServerDatabaseClient, err := databases.NewDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerDatabases client: %+v", err)
	}
	o.Configure(flexibleServerDatabaseClient.Client, o.Authorizers.ResourceManager)

	flexibleServerConfigurationsClient, err := configurations.NewConfigurationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerConfigurations client: %+v", err)
	}
	o.Configure(flexibleServerConfigurationsClient.Client, o.Authorizers.ResourceManager)

	flexibleServerAdministratorsClient, err := administratormicrosoftentras.NewAdministratorMicrosoftEntrasClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building FlexibleServerAdministrators client: %+v", err)
	}
	o.Configure(flexibleServerAdministratorsClient.Client, o.Authorizers.ResourceManager)

	virtualEndpointClient, err := virtualendpoints.NewVirtualEndpointsClientWithBaseURI(o.Environment.ResourceManager)
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
