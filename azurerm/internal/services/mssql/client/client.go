package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	ServerSecurityAlertPoliciesClient                  *sql.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *sql.ServerVulnerabilityAssessmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ElasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ElasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	DatabaseVulnerabilityAssessmentRuleBaselinesClient := sql.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	ServerSecurityAlertPoliciesClient := sql.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	ServerVulnerabilityAssessmentsClient := sql.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ElasticPoolsClient: &ElasticPoolsClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &DatabaseVulnerabilityAssessmentRuleBaselinesClient,
		ServerSecurityAlertPoliciesClient:                  &ServerSecurityAlertPoliciesClient,
		ServerVulnerabilityAssessmentsClient:               &ServerVulnerabilityAssessmentsClient,
	}
}
