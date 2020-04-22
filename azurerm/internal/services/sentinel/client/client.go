package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	AlertRulesClient *securityinsight.AlertRulesClient
	ActionsClient    *securityinsight.ActionsClient
}

func NewClient(o *common.ClientOptions) *Client {
	actionsClient := securityinsight.NewActionsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&actionsClient.Client, o.ResourceManagerAuthorizer)

	alertRulesClient := securityinsight.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ActionsClient:    &actionsClient,
		AlertRulesClient: &alertRulesClient,
	}
}
