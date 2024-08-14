// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/backupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databasesecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/elasticpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/longtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/replicationlinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/restorabledroppeddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverconnectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/transparentdataencryptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupShortTermRetentionPoliciesClient             *backupshorttermretentionpolicies.BackupShortTermRetentionPoliciesClient
	DatabaseExtendedBlobAuditingPoliciesClient         *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	DatabaseSecurityAlertPoliciesClient                *databasesecurityalertpolicies.DatabaseSecurityAlertPoliciesClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	DatabasesClient                                    *databases.DatabasesClient
	ElasticPoolsClient                                 *elasticpools.ElasticPoolsClient
	EncryptionProtectorClient                          *sql.EncryptionProtectorsClient
	FailoverGroupsClient                               *sql.FailoverGroupsClient
	FirewallRulesClient                                *sql.FirewallRulesClient
	GeoBackupPoliciesClient                            *geobackuppolicies.GeoBackupPoliciesClient
	JobAgentsClient                                    *sql.JobAgentsClient
	JobCredentialsClient                               *sql.JobCredentialsClient
	LongTermRetentionPoliciesClient                    *longtermretentionpolicies.LongTermRetentionPoliciesClient
	OutboundFirewallRulesClient                        *sql.OutboundFirewallRulesClient
	ReplicationLinksClient                             *replicationlinks.ReplicationLinksClient
	LegacyReplicationLinksClient                       *sql.ReplicationLinksClient
	RestorableDroppedDatabasesClient                   *restorabledroppeddatabases.RestorableDroppedDatabasesClient
	ServerAzureADAdministratorsClient                  *serverazureadadministrators.ServerAzureADAdministratorsClient
	ServerAzureADOnlyAuthenticationsClient             *serverazureadonlyauthentications.ServerAzureADOnlyAuthenticationsClient
	ServerConnectionPoliciesClient                     *serverconnectionpolicies.ServerConnectionPoliciesClient
	ServerDNSAliasClient                               *sql.ServerDNSAliasesClient
	ServerExtendedBlobAuditingPoliciesClient           *sql.ExtendedServerBlobAuditingPoliciesClient
	ServerDevOpsAuditSettingsClient                    *sql.ServerDevOpsAuditSettingsClient
	ServerKeysClient                                   *sql.ServerKeysClient
	ServerSecurityAlertPoliciesClient                  *serversecurityalertpolicies.ServerSecurityAlertPoliciesClient
	LegacyServerSecurityAlertPoliciesClient            *sql.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *sql.ServerVulnerabilityAssessmentsClient
	ServersClient                                      *servers.ServersClient
	TransparentDataEncryptionsClient                   *transparentdataencryptions.TransparentDataEncryptionsClient
	VirtualMachinesAvailabilityGroupListenersClient    *availabilitygrouplisteners.AvailabilityGroupListenersClient
	VirtualMachinesClient                              *sqlvirtualmachines.SqlVirtualMachinesClient
	VirtualMachineGroupsClient                         *sqlvirtualmachinegroups.SqlVirtualMachineGroupsClient
	VirtualNetworkRulesClient                          *sql.VirtualNetworkRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupShortTermRetentionPoliciesClient, err := backupshorttermretentionpolicies.NewBackupShortTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Backup Short Term Retention Policies Client: %+v", err)
	}
	o.Configure(backupShortTermRetentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	databaseExtendedBlobAuditingPoliciesClient := sql.NewExtendedDatabaseBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseSecurityAlertPoliciesClient, err := databasesecurityalertpolicies.NewDatabaseSecurityAlertPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Databases Security Alert Policies Client: %+v", err)
	}
	o.Configure(databaseSecurityAlertPoliciesClient.Client, o.Authorizers.ResourceManager)

	databaseVulnerabilityAssessmentRuleBaselinesClient := sql.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	databasesClient, err := databases.NewDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Databases Client: %+v", err)
	}
	o.Configure(databasesClient.Client, o.Authorizers.ResourceManager)

	elasticPoolsClient, err := elasticpools.NewElasticPoolsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ElasticPools Client: %+v", err)
	}
	o.Configure(elasticPoolsClient.Client, o.Authorizers.ResourceManager)

	encryptionProtectorClient := sql.NewEncryptionProtectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&encryptionProtectorClient.Client, o.ResourceManagerAuthorizer)

	failoverGroupsClient := sql.NewFailoverGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&failoverGroupsClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient := sql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	geoBackupPoliciesClient, err := geobackuppolicies.NewGeoBackupPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Geo Backup Policies Client: %+v", err)
	}
	o.Configure(geoBackupPoliciesClient.Client, o.Authorizers.ResourceManager)

	jobAgentsClient := sql.NewJobAgentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobAgentsClient.Client, o.ResourceManagerAuthorizer)

	jobCredentialsClient := sql.NewJobCredentialsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobCredentialsClient.Client, o.ResourceManagerAuthorizer)

	longTermRetentionPoliciesClient, err := longtermretentionpolicies.NewLongTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Long Term Retention Policies Client: %+v", err)
	}
	o.Configure(longTermRetentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	outboundFirewallRulesClient := sql.NewOutboundFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&outboundFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	replicationLinksClient, err := replicationlinks.NewReplicationLinksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Replication Links Client: %+v", err)
	}
	o.Configure(replicationLinksClient.Client, o.Authorizers.ResourceManager)

	// NOTE: Remove once Azure Bug 2805551 ReplicationLink API ListByDatabase missed subsubcriptionId in partnerDatabaseId in response body has been released
	legacyReplicationLinksClient := sql.NewReplicationLinksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyReplicationLinksClient.Client, o.ResourceManagerAuthorizer)

	restorableDroppedDatabasesClient, err := restorabledroppeddatabases.NewRestorableDroppedDatabasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Restorable Dropped Databases Client: %+v", err)
	}
	o.Configure(restorableDroppedDatabasesClient.Client, o.Authorizers.ResourceManager)

	serverAzureADAdministratorsClient, err := serverazureadadministrators.NewServerAzureADAdministratorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Azure Active Directory Administrators Client: %+v", err)
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

	legacyServerSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyServerSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient, err := serversecurityalertpolicies.NewServerSecurityAlertPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Security Alert Policies Client: %+v", err)
	}
	o.Configure(serverSecurityAlertPoliciesClient.Client, o.Authorizers.ResourceManager)

	serverVulnerabilityAssessmentsClient := sql.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	serversClient, err := servers.NewServersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Client: %+v", err)
	}
	o.Configure(serversClient.Client, o.Authorizers.ResourceManager)

	transparentDataEncryptionsClient, err := transparentdataencryptions.NewTransparentDataEncryptionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Transparent Data Encryptions Client: %+v", err)
	}
	o.Configure(transparentDataEncryptionsClient.Client, o.Authorizers.ResourceManager)

	virtualMachinesAvailabilityGroupListenersClient, err := availabilitygrouplisteners.NewAvailabilityGroupListenersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Virtual Machines Availability Group Listeners Client: %+v", err)
	}
	o.Configure(virtualMachinesAvailabilityGroupListenersClient.Client, o.Authorizers.ResourceManager)

	virtualMachinesClient, err := sqlvirtualmachines.NewSqlVirtualMachinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Virtual Machines Client: %+v", err)
	}
	o.Configure(virtualMachinesClient.Client, o.Authorizers.ResourceManager)

	virtualMachineGroupsClient, err := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Virtual Machine Groups Client: %+v", err)
	}
	o.Configure(virtualMachineGroupsClient.Client, o.Authorizers.ResourceManager)

	virtualNetworkRulesClient := sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		// Clients using the Track1 SDK which need to be gradually switched over to `hashicorp/go-azure-sdk`
		DatabaseExtendedBlobAuditingPoliciesClient:         &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &databaseVulnerabilityAssessmentRuleBaselinesClient,
		EncryptionProtectorClient:                          &encryptionProtectorClient,
		FailoverGroupsClient:                               &failoverGroupsClient,
		FirewallRulesClient:                                &firewallRulesClient,
		JobAgentsClient:                                    &jobAgentsClient,
		JobCredentialsClient:                               &jobCredentialsClient,
		OutboundFirewallRulesClient:                        &outboundFirewallRulesClient,
		ServerDNSAliasClient:                               &serverDNSAliasClient,
		ServerDevOpsAuditSettingsClient:                    &serverDevOpsAuditSettingsClient,
		ServerExtendedBlobAuditingPoliciesClient:           &serverExtendedBlobAuditingPoliciesClient,
		ServerKeysClient:                                   &serverKeysClient,
		ServerVulnerabilityAssessmentsClient:               &serverVulnerabilityAssessmentsClient,
		VirtualMachinesAvailabilityGroupListenersClient:    virtualMachinesAvailabilityGroupListenersClient,
		VirtualMachinesClient:                              virtualMachinesClient,
		VirtualMachineGroupsClient:                         virtualMachineGroupsClient,
		VirtualNetworkRulesClient:                          &virtualNetworkRulesClient,

		// Legacy Clients
		LegacyServerSecurityAlertPoliciesClient: &legacyServerSecurityAlertPoliciesClient,
		LegacyReplicationLinksClient:            &legacyReplicationLinksClient,

		// 2023-02-01-preview Clients
		BackupShortTermRetentionPoliciesClient: backupShortTermRetentionPoliciesClient,
		DatabasesClient:                        databasesClient,
		DatabaseSecurityAlertPoliciesClient:    databaseSecurityAlertPoliciesClient,
		ElasticPoolsClient:                     elasticPoolsClient,
		GeoBackupPoliciesClient:                geoBackupPoliciesClient,
		LongTermRetentionPoliciesClient:        longTermRetentionPoliciesClient,
		ReplicationLinksClient:                 replicationLinksClient,
		RestorableDroppedDatabasesClient:       restorableDroppedDatabasesClient,
		ServerAzureADAdministratorsClient:      serverAzureADAdministratorsClient,
		ServerAzureADOnlyAuthenticationsClient: serverAzureADOnlyAuthenticationsClient,
		ServerConnectionPoliciesClient:         serverConnectionPoliciesClient,
		ServerSecurityAlertPoliciesClient:      serverSecurityAlertPoliciesClient,
		TransparentDataEncryptionsClient:       transparentDataEncryptionsClient,
		ServersClient:                          serversClient,
	}, nil
}
