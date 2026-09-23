// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertruletemplates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/dataconnectors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/metadata"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/threatintelligence"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/sentinelonboardingstates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlistitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AlertRulesClient         *alertrules.AlertRulesClient
	AlertRuleTemplatesClient *alertruletemplates.AlertRuleTemplatesClient
	AutomationRulesClient    *automationrules.AutomationRulesClient
	DataConnectorsClient     *dataconnectors.DataConnectorsClient
	WatchlistsClient         *watchlists.WatchlistsClient
	WatchlistItemsClient     *watchlistitems.WatchlistItemsClient
	OnboardingStatesClient   *sentinelonboardingstates.SentinelOnboardingStatesClient
	AnalyticsSettingsClient  *securitymlanalyticssettings.SecurityMLAnalyticsSettingsClient
	ThreatIntelligenceClient *threatintelligence.ThreatIntelligenceClient
	MetadataClient           *metadata.MetadataClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	alertRulesClient, err := alertrules.NewAlertRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Alert Rules Client: %+v", err)
	}
	o.Configure(alertRulesClient.Client, o.Authorizers.ResourceManager)

	alertRuleTemplatesClient, err := alertruletemplates.NewAlertRuleTemplatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Alert Rule Templates Client: %+v", err)
	}
	o.Configure(alertRuleTemplatesClient.Client, o.Authorizers.ResourceManager)

	automationRulesClient, err := automationrules.NewAutomationRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Automation Rules Client: %+v", err)
	}
	o.Configure(automationRulesClient.Client, o.Authorizers.ResourceManager)

	dataConnectorsClient, err := dataconnectors.NewDataConnectorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Connectors Client: %+v", err)
	}
	o.Configure(dataConnectorsClient.Client, o.Authorizers.ResourceManager)

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

	analyticsSettingsClient, err := securitymlanalyticssettings.NewSecurityMLAnalyticsSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Security ML Analytics Settings Client: %+v", err)
	}
	o.Configure(analyticsSettingsClient.Client, o.Authorizers.ResourceManager)

	threatIntelligenceClient, err := threatintelligence.NewThreatIntelligenceClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Threat Intelligence Client: %+v", err)
	}
	o.Configure(threatIntelligenceClient.Client, o.Authorizers.ResourceManager)

	metadataClient, err := metadata.NewMetadataClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Metadata Client: %+v", err)
	}
	o.Configure(metadataClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AlertRulesClient:         alertRulesClient,
		AlertRuleTemplatesClient: alertRuleTemplatesClient,
		AutomationRulesClient:    automationRulesClient,
		DataConnectorsClient:     dataConnectorsClient,
		WatchlistsClient:         watchListsClient,
		WatchlistItemsClient:     watchListItemsClient,
		OnboardingStatesClient:   onboardingStatesClient,
		AnalyticsSettingsClient:  analyticsSettingsClient,
		ThreatIntelligenceClient: threatIntelligenceClient,
		MetadataClient:           metadataClient,
	}, nil
}
