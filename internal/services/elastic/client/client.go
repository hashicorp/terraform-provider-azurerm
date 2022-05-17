package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/sdk/2020-07-01/monitorsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/sdk/2020-07-01/rules"
)

type Client struct {
	MonitorClient *monitorsresource.MonitorsResourceClient
	TagRuleClient *rules.RulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := monitorsresource.NewMonitorsResourceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&monitorClient.Client, o.ResourceManagerAuthorizer)

	tagRuleClient := rules.NewRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&tagRuleClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorClient: &monitorClient,
		TagRuleClient: &tagRuleClient,
	}
}
