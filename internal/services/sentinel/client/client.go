package client

import (
	alertruletemplates "github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2021-09-01-preview/securityinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/automationrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type Client struct {
	AlertRulesClient         *alertrules.AlertRulesClient
	AlertRuleTemplatesClient *alertruletemplates.AlertRuleTemplatesClient
	AutomationRulesClient    *automationrules.AutomationRulesClient
	DataConnectorsClient     *securityinsight.DataConnectorsClient
	WatchlistsClient         *securityinsight.WatchlistsClient
	WatchlistItemsClient     *securityinsight.WatchlistItemsClient
	OnboardingStatesClient   *sentinelonboardingstates.SentinelOnboardingStatesClient
	AnalyticsSettingsClient  *securityinsight.SecurityMLAnalyticsSettingsClient
	MetadataClient           *metadata.MetadataClient
}

func NewClient(o *common.ClientOptions) *Client {
	alertRulesClient := alertrules.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&alertRulesClient.Client, o.ResourceManagerAuthorizer)

	alertRuleTemplatesClient := alertruletemplates.NewAlertRuleTemplatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRuleTemplatesClient.Client, o.ResourceManagerAuthorizer)

	automationRulesClient := automationrules.NewAutomationRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&automationRulesClient.Client, o.ResourceManagerAuthorizer)

	dataConnectorsClient := securityinsight.NewDataConnectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dataConnectorsClient.Client, o.ResourceManagerAuthorizer)

	watchListsClient := securityinsight.NewWatchlistsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&watchListsClient.Client, o.ResourceManagerAuthorizer)

	watchListItemsClient := securityinsight.NewWatchlistItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&watchListItemsClient.Client, o.ResourceManagerAuthorizer)

	onboardingStatesClient := sentinelonboardingstates.NewSentinelOnboardingStatesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&onboardingStatesClient.Client, o.ResourceManagerAuthorizer)

	analyticsSettingsClient := securityinsight.NewSecurityMLAnalyticsSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&analyticsSettingsClient.Client, o.ResourceManagerAuthorizer)

	metadataClient := metadata.NewMetadataClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&metadataClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AlertRulesClient:         &alertRulesClient,
		AlertRuleTemplatesClient: &alertRuleTemplatesClient,
		AutomationRulesClient:    &automationRulesClient,
		DataConnectorsClient:     &dataConnectorsClient,
		WatchlistsClient:         &watchListsClient,
		WatchlistItemsClient:     &watchListItemsClient,
		OnboardingStatesClient:   &onboardingStatesClient,
		AnalyticsSettingsClient:  &analyticsSettingsClient,
		MetadataClient:           &metadataClient,
	}
}
