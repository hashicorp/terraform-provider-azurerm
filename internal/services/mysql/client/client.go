// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest"
	flexibleServers_v2021_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2021-05-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/azureadadministrators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FlexibleServers *flexibleServers_v2021_05_01.Client

	ConfigurationsClient              *mysql.ConfigurationsClient
	DatabasesClient                   *mysql.DatabasesClient
	FirewallRulesClient               *mysql.FirewallRulesClient
	ServersClient                     *mysql.ServersClient
	ServerKeysClient                  *mysql.ServerKeysClient
	ServerSecurityAlertPoliciesClient *mysql.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient         *mysql.VirtualNetworkRulesClient
	ServerAdministratorsClient        *mysql.ServerAdministratorsClient

	AzureADAdministratorsClient *azureadadministrators.AzureADAdministratorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	flexibleServersMetaClient := flexibleServers_v2021_05_01.NewClientWithBaseURI(o.ResourceManagerEndpoint, func(c *autorest.Client) {
		o.ConfigureClient(c, o.ResourceManagerAuthorizer)
	})

	azureADAdministratorsClient := azureadadministrators.NewAzureADAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&azureADAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	ConfigurationsClient := mysql.NewConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	DatabasesClient := mysql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabasesClient.Client, o.ResourceManagerAuthorizer)

	FirewallRulesClient := mysql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&FirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerClient := servers.NewServersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerClient.Client, o.ResourceManagerAuthorizer)

	flexibleServerFailoverClient := serverfailover.NewServerFailoverClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&flexibleServerFailoverClient.Client, o.ResourceManagerAuthorizer)

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
		FlexibleServers: &flexibleServersMetaClient,

		// TODO: switch to using the Meta Clients
		AzureADAdministratorsClient:       &azureADAdministratorsClient,
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
