package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2022-07-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2022-07-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorsClient *monitors.MonitorsClient
	TagRulesClient *tagrules.TagRulesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	monitorsClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(monitorsClient.Client, o.Authorizers.ResourceManager)

	tagRulesClient, err := tagrules.NewTagRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, err
	}
	o.Configure(tagRulesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		MonitorsClient: monitorsClient,
		TagRulesClient: tagRulesClient,
	}, nil
}
