package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AlertRulesClient         *securityinsight.AlertRulesClient
	AlertRuleTemplatesClient *securityinsight.AlertRuleTemplatesClient
}

func NewClient(o *common.ClientOptions) *Client {
	alertRulesClient := securityinsight.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRulesClient.Client, o.ResourceManagerAuthorizer)

	alertRuleTemplatesClient := securityinsight.NewAlertRuleTemplatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRuleTemplatesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AlertRulesClient:         &alertRulesClient,
		AlertRuleTemplatesClient: &alertRuleTemplatesClient,
	}
}
