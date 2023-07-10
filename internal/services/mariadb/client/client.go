// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/virtualnetworkrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ConfigurationsClient      *configurations.ConfigurationsClient
	DatabasesClient           *databases.DatabasesClient
	FirewallRulesClient       *firewallrules.FirewallRulesClient
	ServersClient             *servers.ServersClient
	VirtualNetworkRulesClient *virtualnetworkrules.VirtualNetworkRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	configurationsClient := configurations.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&configurationsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := databases.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := firewallrules.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	ServersClient := servers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ServersClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworkRulesClient := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ConfigurationsClient:      &configurationsClient,
		DatabasesClient:           &DatabasesClient,
		FirewallRulesClient:       &FirewallRulesClient,
		ServersClient:             &ServersClient,
		VirtualNetworkRulesClient: &VirtualNetworkRulesClient,
	}
}
