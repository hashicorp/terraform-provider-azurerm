package client

import (
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient  *monitors.MonitorsClient
	TagRulesClient *tagrules.TagRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	monitorClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dynatrace Monitor client: %+v", err)
	}
	o.Configure(monitorClient.Client, o.Authorizers.ResourceManager)

	tagRulesClient, err := tagrules.NewTagRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dynatrace Tag Rules client: %+v", err)
	}
	o.Configure(tagRulesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorClient:  monitorClient,
		TagRulesClient: tagRulesClient,
	}, nil
}
