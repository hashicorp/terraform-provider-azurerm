package datalake

import (
	analytics "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	store "github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	// Data Lake Store
	StoreAccountsClient      store.AccountsClient
	StoreFirewallRulesClient store.FirewallRulesClient
	StoreFilesClient         filesystem.Client

	// Data Lake Analytics
	AnalyticsAccountsClient      analytics.AccountsClient
	AnalyticsFirewallRulesClient analytics.FirewallRulesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.StoreAccountsClient = store.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.StoreAccountsClient.Client, o.ResourceManagerAuthorizer)

	c.StoreFirewallRulesClient = store.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.StoreFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	c.StoreFilesClient = filesystem.NewClient()
	o.ConfigureClient(&c.StoreFilesClient.Client, o.ResourceManagerAuthorizer)

	c.AnalyticsAccountsClient = analytics.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AnalyticsAccountsClient.Client, o.ResourceManagerAuthorizer)

	c.AnalyticsFirewallRulesClient = analytics.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AnalyticsFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
