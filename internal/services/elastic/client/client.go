package client

import (
	"github.com/Azure/azure-sdk-for-go/services/elastic/mgmt/2020-07-01/elastic"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient *elastic.MonitorsClient
	TagRuleClient *elastic.TagRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := elastic.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitorsClient.Client, o.ResourceManagerAuthorizer)

	tagRuleClient := elastic.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorClient: &monitorsClient,
		TagRuleClient: &tagRulesClient,
	}
}
