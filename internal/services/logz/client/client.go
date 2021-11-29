package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient *logz.MonitorsClient
	TagRuleClient *logz.TagRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := logz.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitorClient.Client, o.ResourceManagerAuthorizer)

	tagRuleClient := logz.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagRuleClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorClient: &monitorClient,
		TagRuleClient: &tagRuleClient,
	}
}
