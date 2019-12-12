package client

import (
	sql201703 "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	sql201806 "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *sql201703.DatabaseVulnerabilityAssessmentRuleBaselinesClient
	ServerSecurityAlertPoliciesClient                  *sql201703.ServerSecurityAlertPoliciesClient
	ServerVulnerabilityAssessmentsClient               *sql201806.ServerVulnerabilityAssessmentsClient
}

func NewClient(o *common.ClientOptions) *Client {
	ElasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ElasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	DatabaseVulnerabilityAssessmentRuleBaselinesClient := sql201703.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	ServerSecurityAlertPoliciesClient := sql201703.NewServerSecurityAlertPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerSecurityAlertPoliciesClient.Client, o.ResourceManagerAuthorizer)

	ServerVulnerabilityAssessmentsClient := sql201806.NewServerVulnerabilityAssessmentsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ServerVulnerabilityAssessmentsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ElasticPoolsClient: &ElasticPoolsClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &DatabaseVulnerabilityAssessmentRuleBaselinesClient,
		ServerSecurityAlertPoliciesClient:                  &ServerSecurityAlertPoliciesClient,
		ServerVulnerabilityAssessmentsClient:               &ServerVulnerabilityAssessmentsClient,
	}
}
