package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2021-09-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2021-09-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient  *monitors.MonitorsClient
	TagRulesClient *tagrules.TagRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := monitors.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&monitorClient.Client, o.ResourceManagerAuthorizer)

	tagRulesClient := tagrules.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&tagRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorClient:  &monitorClient,
		TagRulesClient: &tagRulesClient,
	}
}
