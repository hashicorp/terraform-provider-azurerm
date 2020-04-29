package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	DatabasesClient                                    *sql.DatabasesClient
	DatabaseExtendedBlobAuditingPoliciesClient         *sql.ExtendedDatabaseBlobAuditingPoliciesClient
	DatabaseThreatDetectionPoliciesClient              *sql.DatabaseThreatDetectionPoliciesClient
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	ServersClient                                      *sql.ServersClient
	ServerSecurityAlertPoliciesClient                  *sql.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *sql.ServerVulnerabilityAssessmentsClient
	VirtualMachinesClient                              *sqlvirtualmachine.SQLVirtualMachinesClient
}

func NewClient(o *common.ClientOptions) *Client {
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

	serverSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	serversClient := sql.NewServersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serversClient.Client, o.ResourceManagerAuthorizer)

	serverVulnerabilityAssessmentsClient := sql.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&serverVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	sqlVirtualMachinesClient := sqlvirtualmachine.NewSQLVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sqlVirtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		DatabasesClient: &databasesClient,
		DatabaseExtendedBlobAuditingPoliciesClient:         &databaseExtendedBlobAuditingPoliciesClient,
		DatabaseThreatDetectionPoliciesClient:              &databaseThreatDetectionPoliciesClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &databaseVulnerabilityAssessmentRuleBaselinesClient,
		ElasticPoolsClient:                                 &elasticPoolsClient,
		ServersClient:                                      &serversClient,
		ServerSecurityAlertPoliciesClient:                  &serverSecurityAlertPoliciesClient,
		ServerVulnerabilityAssessmentsClient:               &serverVulnerabilityAssessmentsClient,
		VirtualMachinesClient:                              &sqlVirtualMachinesClient,
	}
}
