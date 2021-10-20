package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	msi "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	DatabasesClient                            *sql.DatabasesClient
	DatabaseThreatDetectionPoliciesClient      *sql.DatabaseThreatDetectionPoliciesClient
	ElasticPoolsClient                         *sql.ElasticPoolsClient
	DatabaseExtendedBlobAuditingPoliciesClient *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	FirewallRulesClient                        *sql.FirewallRulesClient
	FailoverGroupsClient                       *sql.FailoverGroupsClient
	ManagedInstancesClient                     *msi.ManagedInstancesClient
	ManagedDatabasesClient                     *msi.ManagedDatabasesClient
	ServersClient                              *sql.ServersClient
	ServerExtendedBlobAuditingPoliciesClient   *sql.ExtendedServerBlobAuditingPoliciesClient
	ServerConnectionPoliciesClient             *sql.ServerConnectionPoliciesClient
	ServerAzureADAdministratorsClient          *sql.ServerAzureADAdministratorsClient
	ServerSecurityAlertPoliciesClient          *sql.ServerSecurityAlertPoliciesClient
	VirtualNetworkRulesClient                  *sql.VirtualNetworkRulesClient
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

	managedInstancesClient := msi.NewManagedInstancesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedInstancesClient.Client, o.ResourceManagerAuthorizer)

	managedDatabasesClient := msi.NewManagedDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&managedDatabasesClient.Client, o.ResourceManagerAuthorizer)

	serversClient := sql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serversClient.Client, o.ResourceManagerAuthorizer)

	serverConnectionPoliciesClient := sql.NewServerConnectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverConnectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverAzureADAdministratorsClient := sql.NewServerAzureADAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAzureADAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	serverExtendedBlobAuditingPoliciesClient := sql.NewExtendedServerBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatabasesClient: &databasesClient,
		DatabaseExtendedBlobAuditingPoliciesClient: &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseThreatDetectionPoliciesClient:      &databaseThreatDetectionPoliciesClient,
		ElasticPoolsClient:                         &elasticPoolsClient,
		FailoverGroupsClient:                       &failoverGroupsClient,
		FirewallRulesClient:                        &firewallRulesClient,
		ManagedInstancesClient:                     &managedInstancesClient,
		ManagedDatabasesClient:                     &managedDatabasesClient,
		ServersClient:                              &serversClient,
		ServerAzureADAdministratorsClient:          &serverAzureADAdministratorsClient,
		ServerConnectionPoliciesClient:             &serverConnectionPoliciesClient,
		ServerExtendedBlobAuditingPoliciesClient:   &serverExtendedBlobAuditingPoliciesClient,
		ServerSecurityAlertPoliciesClient:          &serverSecurityAlertPoliciesClient,
		VirtualNetworkRulesClient:                  &virtualNetworkRulesClient,
	}
}
