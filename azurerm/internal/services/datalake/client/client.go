package client

import (
	analytics "github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	store "github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	// Data Lake Store
	StoreAccountsClient      *store.AccountsClient
	StoreFirewallRulesClient *store.FirewallRulesClient
	StoreFilesClient         *filesystem.Client

	// Data Lake Analytics
	AnalyticsAccountsClient      *analytics.AccountsClient
	AnalyticsFirewallRulesClient *analytics.FirewallRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	StoreAccountsClient := store.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&StoreAccountsClient.Client, o.ResourceManagerAuthorizer)

	StoreFirewallRulesClient := store.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&StoreFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	StoreFilesClient := filesystem.NewClient()
	o.ConfigureClient(&StoreFilesClient.Client, o.ResourceManagerAuthorizer)

	AnalyticsAccountsClient := analytics.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AnalyticsAccountsClient.Client, o.ResourceManagerAuthorizer)

	AnalyticsFirewallRulesClient := analytics.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AnalyticsFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		StoreAccountsClient:          &StoreAccountsClient,
		StoreFirewallRulesClient:     &StoreFirewallRulesClient,
		StoreFilesClient:             &StoreFilesClient,
		AnalyticsAccountsClient:      &AnalyticsAccountsClient,
		AnalyticsFirewallRulesClient: &AnalyticsFirewallRulesClient,
	}
}
