// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	alertruletemplates "github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2021-09-01-preview/securityinsight" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlistitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	securityinsight "github.com/jackofallops/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

type Client struct {
	AlertRulesClient         *alertrules.AlertRulesClient
	AlertRuleTemplatesClient *alertruletemplates.AlertRuleTemplatesClient
	AutomationRulesClient    *automationrules.AutomationRulesClient
	DataConnectorsClient     *securityinsight.DataConnectorsClient
	WatchlistsClient         *watchlists.WatchlistsClient
	WatchlistItemsClient     *watchlistitems.WatchlistItemsClient
	OnboardingStatesClient   *sentinelonboardingstates.SentinelOnboardingStatesClient
	AnalyticsSettingsClient  *securityinsight.SecurityMLAnalyticsSettingsClient
	ThreatIntelligenceClient *securityinsight.ThreatIntelligenceIndicatorClient
	MetadataClient           *metadata.MetadataClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	alertRulesClient, err := alertrules.NewAlertRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Alert Rules Client: %+v", err)
	}
	o.Configure(alertRulesClient.Client, o.Authorizers.ResourceManager)

	alertRuleTemplatesClient := alertruletemplates.NewAlertRuleTemplatesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&alertRuleTemplatesClient.Client, o.ResourceManagerAuthorizer)

	automationRulesClient, err := automationrules.NewAutomationRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Automation Rules Client: %+v", err)
	}
	o.Configure(automationRulesClient.Client, o.Authorizers.ResourceManager)

	dataConnectorsClient := securityinsight.NewDataConnectorsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&dataConnectorsClient.Client, o.ResourceManagerAuthorizer)

	watchListsClient, err := watchlists.NewWatchlistsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Watch Lists Client: %+v", err)
	}
	o.Configure(watchListsClient.Client, o.Authorizers.ResourceManager)

	watchListItemsClient, err := watchlistitems.NewWatchlistItemsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Watch Lists Items Client: %+v", err)
	}
	o.Configure(watchListItemsClient.Client, o.Authorizers.ResourceManager)

	onboardingStatesClient, err := sentinelonboardingstates.NewSentinelOnboardingStatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Onboarding States Client: %+v", err)
	}
	o.Configure(onboardingStatesClient.Client, o.Authorizers.ResourceManager)

	analyticsSettingsClient := securityinsight.NewSecurityMLAnalyticsSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&analyticsSettingsClient.Client, o.ResourceManagerAuthorizer)

	threatIntelligenceClient := securityinsight.NewThreatIntelligenceIndicatorClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&threatIntelligenceClient.Client, o.ResourceManagerAuthorizer)

	metadataClient, err := metadata.NewMetadataClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Metadata Client: %+v", err)
	}
	o.Configure(metadataClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AlertRulesClient:         alertRulesClient,
		AlertRuleTemplatesClient: &alertRuleTemplatesClient,
		AutomationRulesClient:    automationRulesClient,
		DataConnectorsClient:     &dataConnectorsClient,
		WatchlistsClient:         watchListsClient,
		WatchlistItemsClient:     watchListItemsClient,
		OnboardingStatesClient:   onboardingStatesClient,
		AnalyticsSettingsClient:  &analyticsSettingsClient,
		ThreatIntelligenceClient: &threatIntelligenceClient,
		MetadataClient:           metadataClient,
	}, nil
}
