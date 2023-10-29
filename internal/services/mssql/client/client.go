// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"                                       // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/restorabledroppeddatabases"       // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadadministrators"      // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadonlyauthentications" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverconnectionpolicies"         // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers"                          // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupShortTermRetentionPoliciesClient             *sql.BackupShortTermRetentionPoliciesClient
	DatabaseExtendedBlobAuditingPoliciesClient         *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	DatabaseSecurityAlertPoliciesClient                *sql.DatabaseSecurityAlertPoliciesClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	DatabasesClient                                    *sql.DatabasesClient
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	EncryptionProtectorClient                          *sql.EncryptionProtectorsClient
	FailoverGroupsClient                               *sql.FailoverGroupsClient
	FirewallRulesClient                                *sql.FirewallRulesClient
	GeoBackupPoliciesClient                            *sql.GeoBackupPoliciesClient
	JobAgentsClient                                    *sql.JobAgentsClient
	JobCredentialsClient                               *sql.JobCredentialsClient
	LongTermRetentionPoliciesClient                    *sql.LongTermRetentionPoliciesClient
	OutboundFirewallRulesClient                        *sql.OutboundFirewallRulesClient
	ReplicationLinksClient                             *sql.ReplicationLinksClient
	RestorableDroppedDatabasesClient                   *restorabledroppeddatabases.RestorableDroppedDatabasesClient
	ServerAzureADAdministratorsClient                  *serverazureadadministrators.ServerAzureADAdministratorsClient
	ServerAzureADOnlyAuthenticationsClient             *serverazureadonlyauthentications.ServerAzureADOnlyAuthenticationsClient
	ServerConnectionPoliciesClient                     *serverconnectionpolicies.ServerConnectionPoliciesClient
	ServerDNSAliasClient                               *sql.ServerDNSAliasesClient
	ServerExtendedBlobAuditingPoliciesClient           *sql.ExtendedServerBlobAuditingPoliciesClient
	ServerDevOpsAuditSettingsClient                    *sql.ServerDevOpsAuditSettingsClient
	ServerKeysClient                                   *sql.ServerKeysClient
	ServerSecurityAlertPoliciesClient                  *sql.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *sql.ServerVulnerabilityAssessmentsClient
	ServersClient                                      *servers.ServersClient
	TransparentDataEncryptionsClient                   *sql.TransparentDataEncryptionsClient
	VirtualMachinesAvailabilityGroupListenersClient    *availabilitygrouplisteners.AvailabilityGroupListenersClient
	VirtualMachinesClient                              *sqlvirtualmachines.SqlVirtualMachinesClient
	VirtualMachineGroupsClient                         *sqlvirtualmachinegroups.SqlVirtualMachineGroupsClient
	VirtualNetworkRulesClient                          *sql.VirtualNetworkRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupShortTermRetentionPoliciesClient := sql.NewBackupShortTermRetentionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupShortTermRetentionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseExtendedBlobAuditingPoliciesClient := sql.NewExtendedDatabaseBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseSecurityAlertPoliciesClient := sql.NewDatabaseSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseVulnerabilityAssessmentRuleBaselinesClient := sql.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	databasesClient := sql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databasesClient.Client, o.ResourceManagerAuthorizer)

	elasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&elasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	encryptionProtectorClient := sql.NewEncryptionProtectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&encryptionProtectorClient.Client, o.ResourceManagerAuthorizer)

	failoverGroupsClient := sql.NewFailoverGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&failoverGroupsClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient := sql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	geoBackupPoliciesClient := sql.NewGeoBackupPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&geoBackupPoliciesClient.Client, o.ResourceManagerAuthorizer)

	jobAgentsClient := sql.NewJobAgentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobAgentsClient.Client, o.ResourceManagerAuthorizer)

	jobCredentialsClient := sql.NewJobCredentialsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobCredentialsClient.Client, o.ResourceManagerAuthorizer)

	longTermRetentionPoliciesClient := sql.NewLongTermRetentionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&longTermRetentionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	outboundFirewallRulesClient := sql.NewOutboundFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&outboundFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	replicationLinksClient := sql.NewReplicationLinksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&replicationLinksClient.Client, o.ResourceManagerAuthorizer)

	restorableDroppedDatabasesClient, err := restorabledroppeddatabases.NewRestorableDroppedDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Restorable Dropped Databases Client: %+v", err)
	}
	o.Configure(restorableDroppedDatabasesClient.Client, o.Authorizers.ResourceManager)

	serverAzureADAdministratorsClient, err := serverazureadadministrators.NewServerAzureADAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure Active Directory Administrators Client: %+v", err)
	}
	o.Configure(serverAzureADAdministratorsClient.Client, o.Authorizers.ResourceManager)

	serverAzureADOnlyAuthenticationsClient, err := serverazureadonlyauthentications.NewServerAzureADOnlyAuthenticationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Azure Active Directory Only Authentication Client: %+v", err)
	}
	o.Configure(serverAzureADOnlyAuthenticationsClient.Client, o.Authorizers.ResourceManager)

	serverConnectionPoliciesClient, err := serverconnectionpolicies.NewServerConnectionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Connection Policies Client: %+v", err)
	}
	o.Configure(serverConnectionPoliciesClient.Client, o.Authorizers.ResourceManager)

	serverDNSAliasClient := sql.NewServerDNSAliasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverDNSAliasClient.Client, o.ResourceManagerAuthorizer)

	serverExtendedBlobAuditingPoliciesClient := sql.NewExtendedServerBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverDevOpsAuditSettingsClient := sql.NewServerDevOpsAuditSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverDevOpsAuditSettingsClient.Client, o.ResourceManagerAuthorizer)

	serverKeysClient := sql.NewServerKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverKeysClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverVulnerabilityAssessmentsClient := sql.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	serversClient, err := servers.NewServersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Client: %+v", err)
	}
	o.Configure(serversClient.Client, o.Authorizers.ResourceManager)

	transparentDataEncryptionsClient := sql.NewTransparentDataEncryptionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&transparentDataEncryptionsClient.Client, o.ResourceManagerAuthorizer)

	virtualMachinesAvailabilityGroupListenersClient := availabilitygrouplisteners.NewAvailabilityGroupListenersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualMachinesAvailabilityGroupListenersClient.Client, o.ResourceManagerAuthorizer)

	virtualMachinesClient := sqlvirtualmachines.NewSqlVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	virtualMachineGroupsClient := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualMachineGroupsClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BackupShortTermRetentionPoliciesClient:             &backupShortTermRetentionPoliciesClient,
		DatabaseExtendedBlobAuditingPoliciesClient:         &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseSecurityAlertPoliciesClient:                &databaseSecurityAlertPoliciesClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &databaseVulnerabilityAssessmentRuleBaselinesClient,
		DatabasesClient:                                 &databasesClient,
		ElasticPoolsClient:                              &elasticPoolsClient,
		EncryptionProtectorClient:                       &encryptionProtectorClient,
		FailoverGroupsClient:                            &failoverGroupsClient,
		FirewallRulesClient:                             &firewallRulesClient,
		GeoBackupPoliciesClient:                         &geoBackupPoliciesClient,
		JobAgentsClient:                                 &jobAgentsClient,
		JobCredentialsClient:                            &jobCredentialsClient,
		LongTermRetentionPoliciesClient:                 &longTermRetentionPoliciesClient,
		OutboundFirewallRulesClient:                     &outboundFirewallRulesClient,
		ReplicationLinksClient:                          &replicationLinksClient,
		RestorableDroppedDatabasesClient:                restorableDroppedDatabasesClient,
		ServerAzureADAdministratorsClient:               serverAzureADAdministratorsClient,
		ServerAzureADOnlyAuthenticationsClient:          serverAzureADOnlyAuthenticationsClient,
		ServerConnectionPoliciesClient:                  serverConnectionPoliciesClient,
		ServerDNSAliasClient:                            &serverDNSAliasClient,
		ServerDevOpsAuditSettingsClient:                 &serverDevOpsAuditSettingsClient,
		ServerExtendedBlobAuditingPoliciesClient:        &serverExtendedBlobAuditingPoliciesClient,
		ServerKeysClient:                                &serverKeysClient,
		ServerSecurityAlertPoliciesClient:               &serverSecurityAlertPoliciesClient,
		ServerVulnerabilityAssessmentsClient:            &serverVulnerabilityAssessmentsClient,
		ServersClient:                                   serversClient,
		TransparentDataEncryptionsClient:                &transparentDataEncryptionsClient,
		VirtualMachinesAvailabilityGroupListenersClient: &virtualMachinesAvailabilityGroupListenersClient,
		VirtualMachinesClient:                           &virtualMachinesClient,
		VirtualMachineGroupsClient:                      &virtualMachineGroupsClient,
		VirtualNetworkRulesClient:                       &virtualNetworkRulesClient,
	}, nil
}
