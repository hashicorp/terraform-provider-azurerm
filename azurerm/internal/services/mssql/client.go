package mssql

import (
	vulnerabilitySvc "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ElasticPoolsClient                                 *sql.ElasticPoolsClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient *vulnerabilitySvc.DatabaseVulnerabilityAssessmentRuleBaselinesClient
}

func BuildClient(o *common.ClientOptions) *Client {

	ElasticPoolsClient := sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ElasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	DatabaseVulnerabilityAssessmentRuleBaselinesClient := vulnerabilitySvc.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DatabaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ElasticPoolsClient: &ElasticPoolsClient,
		DatabaseVulnerabilityAssessmentRuleBaselinesClient: &DatabaseVulnerabilityAssessmentRuleBaselinesClient,
	}
}
