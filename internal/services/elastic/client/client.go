package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/legacysdk/elastic/mgmt/2020-07-01/elastic"
)

type Client struct {
	MonitorClient *elastic.MonitorsClient
	TagRuleClient *elastic.TagRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := elastic.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitorClient.Client, o.ResourceManagerAuthorizer)

	tagRuleClient := elastic.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagRuleClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorClient: &monitorClient,
		TagRuleClient: &tagRuleClient,
	}
}
