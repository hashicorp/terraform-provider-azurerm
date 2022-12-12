package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient           *logz.MonitorsClient
	TagRuleClient           *logz.TagRulesClient
	SubAccountClient        *logz.SubAccountClient
	SubAccountTagRuleClient *logz.SubAccountTagRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := logz.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitorClient.Client, o.ResourceManagerAuthorizer)

	tagRuleClient := logz.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagRuleClient.Client, o.ResourceManagerAuthorizer)

	subAccountClient := logz.NewSubAccountClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&subAccountClient.Client, o.ResourceManagerAuthorizer)

	subAccountTagRuleClient := logz.NewSubAccountTagRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&subAccountTagRuleClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorClient:           &monitorClient,
		TagRuleClient:           &tagRuleClient,
		SubAccountClient:        &subAccountClient,
		SubAccountTagRuleClient: &subAccountTagRuleClient,
	}
}
