package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	BackupLongTermRetentionPoliciesClient              *sql.BackupLongTermRetentionPoliciesClient
	BackupShortTermRetentionPoliciesClient             *sql.BackupShortTermRetentionPoliciesClient
	DatabaseExtendedBlobAuditingPoliciesClient         *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	DatabaseThreatDetectionPoliciesClient              *sql.DatabaseThreatDetectionPoliciesClient
	DatabasesClient                                    *sql.DatabasesClient
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	FirewallRulesClient                                *sql.FirewallRulesClient
	JobAgentsClient                                    *sql.JobAgentsClient
	JobCredentialsClient                               *sql.JobCredentialsClient
	ReplicationLinksClient                             *sql.ReplicationLinksClient
	RestorableDroppedDatabasesClient                   *sql.RestorableDroppedDatabasesClient
	ServerAzureADAdministratorsClient                  *sql.ServerAzureADAdministratorsClient
	ServerConnectionPoliciesClient                     *sql.ServerConnectionPoliciesClient
	ServerExtendedBlobAuditingPoliciesClient           *sql.ExtendedServerBlobAuditingPoliciesClient
	ServerSecurityAlertPoliciesClient                  *sql.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *sql.ServerVulnerabilityAssessmentsClient
	ServersClient                                      *sql.ServersClient
	VirtualMachinesClient                              *sqlvirtualmachine.SQLVirtualMachinesClient
	VirtualNetworkRulesClient                          *sql.VirtualNetworkRulesClient
	GeoBackupPoliciesClient                            *sql.GeoBackupPoliciesClient
	EncryptionProtectorClient                          *sql.EncryptionProtectorsClient
	ServerKeysClient                                   *sql.ServerKeysClient
}

func NewClient(o *common.ClientOptions) *Client {
	BackupLongTermRetentionPoliciesClient := sql.NewBackupLongTermRetentionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BackupLongTermRetentionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	BackupShortTermRetentionPoliciesClient := sql.NewBackupShortTermRetentionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&BackupShortTermRetentionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databasesClient := sql.NewDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databasesClient.Client, o.ResourceManagerAuthorizer)

	databaseExtendedBlobAuditingPoliciesClient := sql.NewExtendedDatabaseBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseThreatDetectionPoliciesClient := sql.NewDatabaseThreatDetectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseThreatDetectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	databaseVulnerabilityAssessmentRuleBaselinesClient := sql.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&databaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	elasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&elasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	jobAgentsClient := sql.NewJobAgentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobAgentsClient.Client, o.ResourceManagerAuthorizer)

	jobCredentialsClient := sql.NewJobCredentialsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&jobCredentialsClient.Client, o.ResourceManagerAuthorizer)

	firewallRulesClient := sql.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&firewallRulesClient.Client, o.ResourceManagerAuthorizer)

	replicationLinksClient := sql.NewReplicationLinksClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&replicationLinksClient.Client, o.ResourceManagerAuthorizer)

	restorableDroppedDatabasesClient := sql.NewRestorableDroppedDatabasesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&restorableDroppedDatabasesClient.Client, o.ResourceManagerAuthorizer)

	serverSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverExtendedBlobAuditingPoliciesClient := sql.NewExtendedServerBlobAuditingPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverExtendedBlobAuditingPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serverVulnerabilityAssessmentsClient := sql.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	serverAzureADAdministratorsClient := sql.NewServerAzureADAdministratorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverAzureADAdministratorsClient.Client, o.ResourceManagerAuthorizer)

	serversClient := sql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serversClient.Client, o.ResourceManagerAuthorizer)

	serverConnectionPoliciesClient := sql.NewServerConnectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverConnectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	sqlVirtualMachinesClient := sqlvirtualmachine.NewSQLVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlVirtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	sqlVirtualNetworkRulesClient := sql.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlVirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	geoBackupPoliciesClient := sql.NewGeoBackupPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&geoBackupPoliciesClient.Client, o.ResourceManagerAuthorizer)

	sqlEncryptionProtectorClient := sql.NewEncryptionProtectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlEncryptionProtectorClient.Client, o.ResourceManagerAuthorizer)

	sqlServerKeysClient := sql.NewServerKeysClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlServerKeysClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		BackupLongTermRetentionPoliciesClient:              &BackupLongTermRetentionPoliciesClient,
		BackupShortTermRetentionPoliciesClient:             &BackupShortTermRetentionPoliciesClient,
		DatabasesClient:                                    &databasesClient,
		DatabaseExtendedBlobAuditingPoliciesClient:         &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseThreatDetectionPoliciesClient:              &databaseThreatDetectionPoliciesClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &databaseVulnerabilityAssessmentRuleBaselinesClient,
		EncryptionProtectorClient:                          &sqlEncryptionProtectorClient,
		ElasticPoolsClient:                                 &elasticPoolsClient,
		JobAgentsClient:                                    &jobAgentsClient,
		JobCredentialsClient:                               &jobCredentialsClient,
		FirewallRulesClient:                                &firewallRulesClient,
		ReplicationLinksClient:                             &replicationLinksClient,
		RestorableDroppedDatabasesClient:                   &restorableDroppedDatabasesClient,
		ServerAzureADAdministratorsClient:                  &serverAzureADAdministratorsClient,
		ServersClient:                                      &serversClient,
		ServerExtendedBlobAuditingPoliciesClient:           &serverExtendedBlobAuditingPoliciesClient,
		ServerConnectionPoliciesClient:                     &serverConnectionPoliciesClient,
		ServerSecurityAlertPoliciesClient:                  &serverSecurityAlertPoliciesClient,
		ServerVulnerabilityAssessmentsClient:               &serverVulnerabilityAssessmentsClient,
		VirtualMachinesClient:                              &sqlVirtualMachinesClient,
		VirtualNetworkRulesClient:                          &sqlVirtualNetworkRulesClient,
		GeoBackupPoliciesClient:                            &geoBackupPoliciesClient,
		ServerKeysClient:                                   &sqlServerKeysClient,
	}
}
