package client

import (
	"github.com/Azure/azure-sdk-for-go/services/datalake/store/2016-11-01/filesystem"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	analyticsaccount "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakeanalytics/2016-11-01/accounts"
	analyticsfirewallrules "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakeanalytics/2016-11-01/firewallrules"
	storeaccount "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakestore/2016-11-01/accounts"
	storefirewallrules "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakestore/2016-11-01/firewallrules"
	storevirtualnetworkrules "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/sdk/datalakestore/2016-11-01/virtualnetworkrules"
)

type Client struct {
	// Data Lake Store
	StoreAccountsClient       *storeaccount.AccountsClient
	StoreFirewallRulesClient  *storefirewallrules.FirewallRulesClient
	VirtualNetworkRulesClient *storevirtualnetworkrules.VirtualNetworkRulesClient
	StoreFilesClient          *filesystem.Client

	// Data Lake Analytics
	AnalyticsAccountsClient      *analyticsaccount.AccountsClient
	AnalyticsFirewallRulesClient *analyticsfirewallrules.FirewallRulesClient

	SubscriptionId string
}

func NewClient(o *common.ClientOptions) *Client {
	StoreAccountsClient := storeaccount.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&StoreAccountsClient.Client, o.ResourceManagerAuthorizer)

	StoreFirewallRulesClient := storefirewallrules.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&StoreFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	VirtualNetworkRulesClient := storevirtualnetworkrules.NewVirtualNetworkRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&VirtualNetworkRulesClient.Client, o.ResourceManagerAuthorizer)

	StoreFilesClient := filesystem.NewClient()
	o.ConfigureClient(&StoreFilesClient.Client, o.ResourceManagerAuthorizer)

	AnalyticsAccountsClient := analyticsaccount.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AnalyticsAccountsClient.Client, o.ResourceManagerAuthorizer)

	AnalyticsFirewallRulesClient := analyticsfirewallrules.NewFirewallRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AnalyticsFirewallRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		StoreAccountsClient:          &StoreAccountsClient,
		StoreFirewallRulesClient:     &StoreFirewallRulesClient,
		VirtualNetworkRulesClient:    &VirtualNetworkRulesClient,
		StoreFilesClient:             &StoreFilesClient,
		AnalyticsAccountsClient:      &AnalyticsAccountsClient,
		AnalyticsFirewallRulesClient: &AnalyticsFirewallRulesClient,
		SubscriptionId:               o.SubscriptionId,
	}
}
