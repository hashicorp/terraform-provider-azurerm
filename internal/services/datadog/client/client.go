package client

import (
	"github.com/Azure/azure-sdk-for-go/services/datadog/mgmt/2021-03-01/datadog"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	MonitorsClient                   *datadog.MonitorsClient
	TagRulesClient                   *datadog.TagRulesClient
	SingleSignOnConfigurationsClient *datadog.SingleSignOnConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	monitorsClient := datadog.NewMonitorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&monitorsClient.Client, o.ResourceManagerAuthorizer)

	tagRulesClient := datadog.NewTagRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&tagRulesClient.Client, o.ResourceManagerAuthorizer)

	singleSignOnConfigurationsClient := datadog.NewSingleSignOnConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&singleSignOnConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		MonitorsClient:                   &monitorsClient,
		TagRulesClient:                   &tagRulesClient,
		SingleSignOnConfigurationsClient: &singleSignOnConfigurationsClient,
	}
}
