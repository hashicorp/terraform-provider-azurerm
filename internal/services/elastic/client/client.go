package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2020-07-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2020-07-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
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
