package client

import (
	"github.com/Azure/azure-sdk-for-go/services/logz/mgmt/2020-10-01/logz" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/subaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorClient           *monitors.MonitorsClient
	TagRuleClient           *tagrules.TagRulesClient
	SubAccountClient        *subaccount.SubAccountClient
	SubAccountTagRuleClient *logz.SubAccountTagRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorClient := monitors.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&monitorClient.Client, o.ResourceManagerAuthorizer)

	tagRuleClient := tagrules.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&tagRuleClient.Client, o.ResourceManagerAuthorizer)

	subAccountClient := subaccount.NewSubAccountClientWithBaseURI(o.ResourceManagerEndpoint)
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
