package mssql

import (
	vulnerabilitySvc "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-10-01-preview/sql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	ElasticPoolsClient                                 sql.ElasticPoolsClient
	DatabaseVulnerabilityAssessmentRuleBaselinesClient vulnerabilitySvc.DatabaseVulnerabilityAssessmentRuleBaselinesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.ElasticPoolsClient = sql.NewElasticPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ElasticPoolsClient.Client, o.ResourceManagerAuthorizer)

	c.DatabaseVulnerabilityAssessmentRuleBaselinesClient = vulnerabilitySvc.NewDatabaseVulnerabilityAssessmentRuleBaselinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DatabaseVulnerabilityAssessmentRuleBaselinesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
