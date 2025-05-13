// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/backupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/blobauditing"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databasesecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databasevulnerabilityassessmentrulebaselines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/encryptionprotectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/failovergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/firewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobagents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobcredentials"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobsteps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/jobtargetgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/longtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/outboundfirewallrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/replicationlinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/restorabledroppeddatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverazureadadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverconnectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverdevopsaudit"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverdnsaliases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serverkeys"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servervulnerabilityassessments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/sqlvulnerabilityassessmentssettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/transparentdataencryptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/virtualnetworkrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupShortTermRetentionPoliciesClient             *backupshorttermretentionpolicies.BackupShortTermRetentionPoliciesClient
	BlobAuditingPoliciesClient                         *blobauditing.BlobAuditingClient
	DatabaseSecurityAlertPoliciesClient                *databasesecurityalertpolicies.DatabaseSecurityAlertPoliciesClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *databasevulnerabilityassessmentrulebaselines.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	DatabasesClient                                    *databases.DatabasesClient
	ElasticPoolsClient                                 *elasticpools.ElasticPoolsClient
	EncryptionProtectorClient                          *encryptionprotectors.EncryptionProtectorsClient
	FailoverGroupsClient                               *failovergroups.FailoverGroupsClient
	FirewallRulesClient                                *firewallrules.FirewallRulesClient
	GeoBackupPoliciesClient                            *geobackuppolicies.GeoBackupPoliciesClient
	JobAgentsClient                                    *jobagents.JobAgentsClient
	JobCredentialsClient                               *jobcredentials.JobCredentialsClient
	JobsClient                                         *jobs.JobsClient
	JobStepsClient                                     *jobsteps.JobStepsClient
	JobTargetGroupsClient                              *jobtargetgroups.JobTargetGroupsClient
	LongTermRetentionPoliciesClient                    *longtermretentionpolicies.LongTermRetentionPoliciesClient
	OutboundFirewallRulesClient                        *outboundfirewallrules.OutboundFirewallRulesClient
	ReplicationLinksClient                             *replicationlinks.ReplicationLinksClient
	RestorableDroppedDatabasesClient                   *restorabledroppeddatabases.RestorableDroppedDatabasesClient
	ServerAzureADAdministratorsClient                  *serverazureadadministrators.ServerAzureADAdministratorsClient
	ServerAzureADOnlyAuthenticationsClient             *serverazureadonlyauthentications.ServerAzureADOnlyAuthenticationsClient
	ServerConnectionPoliciesClient                     *serverconnectionpolicies.ServerConnectionPoliciesClient
	ServerDNSAliasClient                               *serverdnsaliases.ServerDnsAliasesClient
	ServerDevOpsAuditSettingsClient                    *serverdevopsaudit.ServerDevOpsAuditClient
	ServerKeysClient                                   *serverkeys.ServerKeysClient
	ServerSecurityAlertPoliciesClient                  *serversecurityalertpolicies.ServerSecurityAlertPoliciesClient
	LegacyServerSecurityAlertPoliciesClient            *sql.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *servervulnerabilityassessments.ServerVulnerabilityAssessmentsClient
	SqlVulnerabilityAssessmentSettingsClient           *sqlvulnerabilityassessmentssettings.SqlVulnerabilityAssessmentsSettingsClient
	ServersClient                                      *servers.ServersClient
	TransparentDataEncryptionsClient                   *transparentdataencryptions.TransparentDataEncryptionsClient
	VirtualMachinesAvailabilityGroupListenersClient    *availabilitygrouplisteners.AvailabilityGroupListenersClient
	VirtualMachinesClient                              *sqlvirtualmachines.SqlVirtualMachinesClient
	VirtualMachineGroupsClient                         *sqlvirtualmachinegroups.SqlVirtualMachineGroupsClient
	VirtualNetworkRulesClient                          *virtualnetworkrules.VirtualNetworkRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	backupShortTermRetentionPoliciesClient, err := backupshorttermretentionpolicies.NewBackupShortTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Backup Short Term Retention Policies Client: %+v", err)
	}
	o.Configure(backupShortTermRetentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	databaseExtendedBlobAuditingPoliciesClient, err := blobauditing.NewBlobAuditingClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Blob Auditing Policies Client: %+v", err)
	}
	o.Configure(databaseExtendedBlobAuditingPoliciesClient.Client, o.Authorizers.ResourceManager)

	databaseSecurityAlertPoliciesClient, err := databasesecurityalertpolicies.NewDatabaseSecurityAlertPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Databases Security Alert Policies Client: %+v", err)
	}
	o.Configure(databaseSecurityAlertPoliciesClient.Client, o.Authorizers.ResourceManager)

	databaseVulnerabilityAssessmentRuleBaselinesClient, err := databasevulnerabilityassessmentrulebaselines.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Database Vulnerability Assessment Rule Baselines Client: %+v", err)
	}
	o.Configure(databaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.Authorizers.ResourceManager)

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

	encryptionProtectorClient, err := encryptionprotectors.NewEncryptionProtectorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Encryption Protectors Client: %+v", err)
	}
	o.Configure(encryptionProtectorClient.Client, o.Authorizers.ResourceManager)

	failoverGroupsClient, err := failovergroups.NewFailoverGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Failover Groups Client: %+v", err)
	}
	o.Configure(failoverGroupsClient.Client, o.Authorizers.ResourceManager)

	firewallRulesClient, err := firewallrules.NewFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Firewall Rules Client: %+v", err)
	}
	o.Configure(firewallRulesClient.Client, o.Authorizers.ResourceManager)

	geoBackupPoliciesClient, err := geobackuppolicies.NewGeoBackupPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Geo Backup Policies Client: %+v", err)
	}
	o.Configure(geoBackupPoliciesClient.Client, o.Authorizers.ResourceManager)

	jobAgentsClient, err := jobagents.NewJobAgentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Job Agents Client: %+v", err)
	}
	o.Configure(jobAgentsClient.Client, o.Authorizers.ResourceManager)

	jobCredentialsClient, err := jobcredentials.NewJobCredentialsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Job Credentials Client: %+v", err)
	}
	o.Configure(jobCredentialsClient.Client, o.Authorizers.ResourceManager)

	jobsClient, err := jobs.NewJobsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Jobs Client: %+v", err)
	}
	o.Configure(jobsClient.Client, o.Authorizers.ResourceManager)

	jobStepsClient, err := jobsteps.NewJobStepsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Job Steps Client: %+v", err)
	}
	o.Configure(jobStepsClient.Client, o.Authorizers.ResourceManager)

	jobTargetGroupsClient, err := jobtargetgroups.NewJobTargetGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Job Target Groups Client: %+v", err)
	}
	o.Configure(jobTargetGroupsClient.Client, o.Authorizers.ResourceManager)

	longTermRetentionPoliciesClient, err := longtermretentionpolicies.NewLongTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Long Term Retention Policies Client: %+v", err)
	}
	o.Configure(longTermRetentionPoliciesClient.Client, o.Authorizers.ResourceManager)

	outboundFirewallRulesClient, err := outboundfirewallrules.NewOutboundFirewallRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Outbound Firewall Rules Client: %+v", err)
	}
	o.Configure(outboundFirewallRulesClient.Client, o.Authorizers.ResourceManager)

	replicationLinksClient, err := replicationlinks.NewReplicationLinksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Replication Links Client: %+v", err)
	}
	o.Configure(replicationLinksClient.Client, o.Authorizers.ResourceManager)

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

	serverDNSAliasClient, err := serverdnsaliases.NewServerDnsAliasesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building server DNS Aliases Client: %+v", err)
	}
	o.Configure(serverDNSAliasClient.Client, o.Authorizers.ResourceManager)

	serverDevOpsAuditSettingsClient, err := serverdevopsaudit.NewServerDevOpsAuditClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building server DevOps Audit Settings Client: %+v", err)
	}
	o.Configure(serverDevOpsAuditSettingsClient.Client, o.Authorizers.ResourceManager)

	serverKeysClient, err := serverkeys.NewServerKeysClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Keys Client: %+v", err)
	}
	o.Configure(serverKeysClient.Client, o.Authorizers.ResourceManager)

	legacyServerSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&legacyServerSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient, err := serversecurityalertpolicies.NewServerSecurityAlertPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Security Alert Policies Client: %+v", err)
	}
	o.Configure(serverSecurityAlertPoliciesClient.Client, o.Authorizers.ResourceManager)

	serverVulnerabilityAssessmentsClient, err := servervulnerabilityassessments.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Server Vulnerability Assessments Client: %+v", err)
	}
	o.Configure(serverVulnerabilityAssessmentsClient.Client, o.Authorizers.ResourceManager)

	sqlVulnerabilityAssessmentsSettingsClient, err := sqlvulnerabilityassessmentssettings.NewSqlVulnerabilityAssessmentsSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SQL Vulnerability Assessments Settings Client: %+v", err)
	}
	o.Configure(sqlVulnerabilityAssessmentsSettingsClient.Client, o.Authorizers.ResourceManager)

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

	virtualNetworkRulesClient, err := virtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Virtual Network Rules Client: %+v", err)
	}
	o.Configure(virtualNetworkRulesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		// Clients using the Track1 SDK which need to be gradually switched over to `hashicorp/go-azure-sdk`
		BlobAuditingPoliciesClient:                         databaseExtendedBlobAuditingPoliciesClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: databaseVulnerabilityAssessmentRuleBaselinesClient,
		EncryptionProtectorClient:                          encryptionProtectorClient,
		FailoverGroupsClient:                               failoverGroupsClient,
		FirewallRulesClient:                                firewallRulesClient,
		JobAgentsClient:                                    jobAgentsClient,
		JobCredentialsClient:                               jobCredentialsClient,
		OutboundFirewallRulesClient:                        outboundFirewallRulesClient,
		ServerDNSAliasClient:                               serverDNSAliasClient,
		ServerDevOpsAuditSettingsClient:                    serverDevOpsAuditSettingsClient,
		ServerKeysClient:                                   serverKeysClient,
		ServerVulnerabilityAssessmentsClient:               serverVulnerabilityAssessmentsClient,
		VirtualMachinesAvailabilityGroupListenersClient:    virtualMachinesAvailabilityGroupListenersClient,
		VirtualMachinesClient:                              virtualMachinesClient,
		VirtualMachineGroupsClient:                         virtualMachineGroupsClient,
		VirtualNetworkRulesClient:                          virtualNetworkRulesClient,

		// Legacy Clients
		LegacyServerSecurityAlertPoliciesClient: &legacyServerSecurityAlertPoliciesClient,

		// 2023-08-01-preview Clients
		BackupShortTermRetentionPoliciesClient:   backupShortTermRetentionPoliciesClient,
		DatabasesClient:                          databasesClient,
		DatabaseSecurityAlertPoliciesClient:      databaseSecurityAlertPoliciesClient,
		ElasticPoolsClient:                       elasticPoolsClient,
		GeoBackupPoliciesClient:                  geoBackupPoliciesClient,
		JobsClient:                               jobsClient,
		JobStepsClient:                           jobStepsClient,
		JobTargetGroupsClient:                    jobTargetGroupsClient,
		LongTermRetentionPoliciesClient:          longTermRetentionPoliciesClient,
		ReplicationLinksClient:                   replicationLinksClient,
		RestorableDroppedDatabasesClient:         restorableDroppedDatabasesClient,
		ServerAzureADAdministratorsClient:        serverAzureADAdministratorsClient,
		ServerAzureADOnlyAuthenticationsClient:   serverAzureADOnlyAuthenticationsClient,
		ServerConnectionPoliciesClient:           serverConnectionPoliciesClient,
		ServerSecurityAlertPoliciesClient:        serverSecurityAlertPoliciesClient,
		SqlVulnerabilityAssessmentSettingsClient: sqlVulnerabilityAssessmentsSettingsClient,
		TransparentDataEncryptionsClient:         transparentDataEncryptionsClient,
		ServersClient:                            serversClient,
	}, nil
}
