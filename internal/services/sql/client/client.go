// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"     // nolint: staticcheck
	msi "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql" // nolint: staticcheck
	sqlv5 "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"             // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DatabasesClient                                 *sql.DatabasesClient
	DatabaseThreatDetectionPoliciesClient           *sql.DatabaseThreatDetectionPoliciesClient
	ElasticPoolsClient                              *sql.ElasticPoolsClient
	DatabaseExtendedBlobAuditingPoliciesClient      *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	FirewallRulesClient                             *sql.FirewallRulesClient
	FailoverGroupsClient                            *sql.FailoverGroupsClient
	InstanceFailoverGroupsClient                    *sqlv5.InstanceFailoverGroupsClient
	ManagedInstancesClient                          *sqlv5.ManagedInstancesClient
	ManagedInstanceAdministratorsClient             *sqlv5.ManagedInstanceAdministratorsClient
	ManagedInstanceAzureADOnlyAuthenticationsClient *sqlv5.ManagedInstanceAzureADOnlyAuthenticationsClient
	ManagedDatabasesClient                          *msi.ManagedDatabasesClient
	ServersClient                                   *sql.ServersClient
	ServerExtendedBlobAuditingPoliciesClient        *sql.ExtendedServerBlobAuditingPoliciesClient
	ServerConnectionPoliciesClient                  *sql.ServerConnectionPoliciesClient
	ServerAzureADAdministratorsClient               *sqlv5.ServerAzureADAdministratorsClient
	ServerAzureADOnlyAuthenticationsClient          *sqlv5.ServerAzureADOnlyAuthenticationsClient
	ServerSecurityAlertPoliciesClient               *sql.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient                       *sql.VirtualNetworkRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	// SQL Azure
	databasesClient := sql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databasesClient.Client, o.ResourceManagerAuthorizer)

	databaseExtendedBlobAuditingPoliciesClient := sql.NewExtendedDatabaseBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseThreatDetectionPoliciesClient := sql.NewDatabaseThreatDetectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseThreatDetectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	elasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&elasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	failoverGroupsClient := sql.NewFailoverGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&failoverGroupsClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient := sql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	instanceFailoverGroupsClient := sqlv5.NewInstanceFailoverGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&instanceFailoverGroupsClient.Client, o.ResourceManagerAuthorizer)

	managedInstancesClient := sqlv5.NewManagedInstancesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstancesClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceAdministratorsClient := sqlv5.NewManagedInstanceAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	managedInstanceAzureADOnlyAuthenticationsClient := sqlv5.NewManagedInstanceAzureADOnlyAuthenticationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstanceAzureADOnlyAuthenticationsClient.Client, o.ResourceManagerAuthorizer)

	managedDatabasesClient := msi.NewManagedDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedDatabasesClient.Client, o.ResourceManagerAuthorizer)

	serversClient := sql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serversClient.Client, o.ResourceManagerAuthorizer)

	serverConnectionPoliciesClient := sql.NewServerConnectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverConnectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverAzureADAdministratorsClient := sqlv5.NewServerAzureADAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAzureADAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	serverAzureADOnlyAuthenticationsClient := sqlv5.NewServerAzureADOnlyAuthenticationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAzureADOnlyAuthenticationsClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	serverExtendedBlobAuditingPoliciesClient := sql.NewExtendedServerBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatabasesClient: &databasesClient,
		DatabaseExtendedBlobAuditingPoliciesClient:      &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseThreatDetectionPoliciesClient:           &databaseThreatDetectionPoliciesClient,
		ElasticPoolsClient:                              &elasticPoolsClient,
		FailoverGroupsClient:                            &failoverGroupsClient,
		FirewallRulesClient:                             &firewallRulesClient,
		InstanceFailoverGroupsClient:                    &instanceFailoverGroupsClient,
		ManagedInstancesClient:                          &managedInstancesClient,
		ManagedInstanceAdministratorsClient:             &managedInstanceAdministratorsClient,
		ManagedInstanceAzureADOnlyAuthenticationsClient: &managedInstanceAzureADOnlyAuthenticationsClient,
		ManagedDatabasesClient:                          &managedDatabasesClient,
		ServersClient:                                   &serversClient,
		ServerAzureADAdministratorsClient:               &serverAzureADAdministratorsClient,
		ServerAzureADOnlyAuthenticationsClient:          &serverAzureADOnlyAuthenticationsClient,
		ServerConnectionPoliciesClient:                  &serverConnectionPoliciesClient,
		ServerExtendedBlobAuditingPoliciesClient:        &serverExtendedBlobAuditingPoliciesClient,
		ServerSecurityAlertPoliciesClient:               &serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:                       &virtualNetworkRulesClient,
	}
}
