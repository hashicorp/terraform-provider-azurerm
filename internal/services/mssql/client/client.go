// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"                                       // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/backupshorttermretentionpolicies" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases"                        // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databasesecurityalertpolicies"    // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/geobackuppolicies"                // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/longtermretentionbackups"         // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/replicationlinks"                 // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/restorabledroppeddatabases"       // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadadministrators"      // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverazureadonlyauthentications" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serverconnectionpolicies"         // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers"                          // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serversecurityalertpolicies"      // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/transparentdataencryptions"       // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BackupShortTermRetentionPoliciesClient             *backupshorttermretentionpolicies.LongTermRetentionPoliciesClient
	DatabaseExtendedBlobAuditingPoliciesClient         *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	DatabaseSecurityAlertPoliciesClient                *databasesecurityalertpolicies.DatabaseSecurityAlertPoliciesClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	DatabasesClient                                    *databases.DatabasesClient
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	EncryptionProtectorClient                          *sql.EncryptionProtectorsClient
	FailoverGroupsClient                               *sql.FailoverGroupsClient
	FirewallRulesClient                                *sql.FirewallRulesClient
	GeoBackupPoliciesClient                            *geobackuppolicies.GeoBackupPoliciesClient
	JobAgentsClient                                    *sql.JobAgentsClient
	JobCredentialsClient                               *sql.JobCredentialsClient
	LongTermRetentionPoliciesClient                    *longtermretentionbackups.LongTermRetentionPoliciesClient
	OutboundFirewallRulesClient                        *sql.OutboundFirewallRulesClient
	ReplicationLinksClient                             *replicationlinks.ReplicationLinksClient
	RestorableDroppedDatabasesClient                   *restorabledroppeddatabases.RestorableDroppedDatabasesClient
	ServerAzureADAdministratorsClient                  *serverazureadadministrators.ServerAzureADAdministratorsClient
	ServerAzureADOnlyAuthenticationsClient             *serverazureadonlyauthentications.ServerAzureADOnlyAuthenticationsClient
	ServerConnectionPoliciesClient                     *serverconnectionpolicies.ServerConnectionPoliciesClient
	ServerDNSAliasClient                               *sql.ServerDNSAliasesClient
	ServerExtendedBlobAuditingPoliciesClient           *sql.ExtendedServerBlobAuditingPoliciesClient
	ServerDevOpsAuditSettingsClient                    *sql.ServerDevOpsAuditSettingsClient
	ServerKeysClient                                   *sql.ServerKeysClient
	ServerSecurityAlertPoliciesClient                  *serversecurityalertpolicies.ServerSecurityAlertPoliciesClient
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

	elasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&elasticPoolsClient.Client, o.ResourceManagerAuthorizer)

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

	longTermRetentionPoliciesClient, err := longtermretentionbackups.NewLongTermRetentionPoliciesClientWithBaseURI(o.Environment.ResourceManager)
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

	virtualMachinesAvailabilityGroupListenersClient := availabilitygrouplisteners.NewAvailabilityGroupListenersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualMachinesAvailabilityGroupListenersClient.Client, o.ResourceManagerAuthorizer)

	virtualMachinesClient := sqlvirtualmachines.NewSqlVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	virtualMachineGroupsClient := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&virtualMachineGroupsClient.Client, o.ResourceManagerAuthorizer)

	virtualNetworkRulesClient := sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&virtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatabaseExtendedBlobAuditingPoliciesClient:         &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &databaseVulnerabilityAssessmentRuleBaselinesClient,
		ElasticPoolsClient:                              &elasticPoolsClient,
		EncryptionProtectorClient:                       &encryptionProtectorClient,
		FailoverGroupsClient:                            &failoverGroupsClient,
		FirewallRulesClient:                             &firewallRulesClient,
		JobAgentsClient:                                 &jobAgentsClient,
		JobCredentialsClient:                            &jobCredentialsClient,
		OutboundFirewallRulesClient:                     &outboundFirewallRulesClient,
		ServerDNSAliasClient:                            &serverDNSAliasClient,
		ServerDevOpsAuditSettingsClient:                 &serverDevOpsAuditSettingsClient,
		ServerExtendedBlobAuditingPoliciesClient:        &serverExtendedBlobAuditingPoliciesClient,
		ServerKeysClient:                                &serverKeysClient,
		ServerVulnerabilityAssessmentsClient:            &serverVulnerabilityAssessmentsClient,
		VirtualMachinesAvailabilityGroupListenersClient: &virtualMachinesAvailabilityGroupListenersClient,
		VirtualMachinesClient:                           &virtualMachinesClient,
		VirtualMachineGroupsClient:                      &virtualMachineGroupsClient,
		VirtualNetworkRulesClient:                       &virtualNetworkRulesClient,

		// 2023-02-01-preview Clients
		BackupShortTermRetentionPoliciesClient: backupShortTermRetentionPoliciesClient,
		DatabasesClient:                        databasesClient,
		DatabaseSecurityAlertPoliciesClient:    databaseSecurityAlertPoliciesClient,
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
