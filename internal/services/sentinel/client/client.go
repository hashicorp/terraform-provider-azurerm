package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AlertRulesClient         *securityinsight.AlertRulesClient
	AlertRuleTemplatesClient *securityinsight.AlertRuleTemplatesClient
	AutomationRulesClient    *securityinsight.AutomationRulesClient
	DataConnectorsClient     *securityinsight.DataConnectorsClient
}

func NewClient(o *common.ClientOptions) *Client {
	alertRulesClient := securityinsight.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRulesClient.Client, o.ResourceManagerAuthorizer)

	alertRuleTemplatesClient := securityinsight.NewAlertRuleTemplatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRuleTemplatesClient.Client, o.ResourceManagerAuthorizer)

	automationRulesClient := securityinsight.NewAutomationRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&automationRulesClient.Client, o.ResourceManagerAuthorizer)

	dataConnectorsClient := securityinsight.NewDataConnectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dataConnectorsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AlertRulesClient:         &alertRulesClient,
		AlertRuleTemplatesClient: &alertRuleTemplatesClient,
		AutomationRulesClient:    &automationRulesClient,
		DataConnectorsClient:     &dataConnectorsClient,
	}
}
