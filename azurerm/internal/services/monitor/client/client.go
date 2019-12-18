package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	// Autoscale Settings
	AutoscaleSettingsClient *insights.AutoscaleSettingsClient

	// Monitor
	ActionGroupsClient               *insights.ActionGroupsClient
	ActivityLogAlertsClient          *insights.ActivityLogAlertsClient
	AlertRulesClient                 *insights.AlertRulesClient
	DiagnosticSettingsClient         *insights.DiagnosticSettingsClient
	DiagnosticSettingsCategoryClient *insights.DiagnosticSettingsCategoryClient
	LogProfilesClient                *insights.LogProfilesClient
	MetricAlertsClient               *insights.MetricAlertsClient
}

func NewClient(o *common.ClientOptions) *Client {
	AutoscaleSettingsClient := insights.NewAutoscaleSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AutoscaleSettingsClient.Client, o.ResourceManagerAuthorizer)

	ActionGroupsClient := insights.NewActionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActionGroupsClient.Client, o.ResourceManagerAuthorizer)

	ActivityLogAlertsClient := insights.NewActivityLogAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActivityLogAlertsClient.Client, o.ResourceManagerAuthorizer)

	AlertRulesClient := insights.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AlertRulesClient.Client, o.ResourceManagerAuthorizer)

	DiagnosticSettingsClient := insights.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DiagnosticSettingsClient.Client, o.ResourceManagerAuthorizer)

	DiagnosticSettingsCategoryClient := insights.NewDiagnosticSettingsCategoryClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DiagnosticSettingsCategoryClient.Client, o.ResourceManagerAuthorizer)

	LogProfilesClient := insights.NewLogProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LogProfilesClient.Client, o.ResourceManagerAuthorizer)

	MetricAlertsClient := insights.NewMetricAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MetricAlertsClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AutoscaleSettingsClient:          &AutoscaleSettingsClient,
		ActionGroupsClient:               &ActionGroupsClient,
		ActivityLogAlertsClient:          &ActivityLogAlertsClient,
		AlertRulesClient:                 &AlertRulesClient,
		DiagnosticSettingsClient:         &DiagnosticSettingsClient,
		DiagnosticSettingsCategoryClient: &DiagnosticSettingsCategoryClient,
		LogProfilesClient:                &LogProfilesClient,
		MetricAlertsClient:               &MetricAlertsClient,
	}
}
