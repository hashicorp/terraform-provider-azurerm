package client

import (
	"github.com/Azure/azure-sdk-for-go/services/aad/mgmt/2017-04-01/aad"
	"github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2020-10-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement"
	classic "github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	// AAD
	AADDiagnosticSettingsClient *aad.DiagnosticSettingsClient

	// Autoscale Settings
	AutoscaleSettingsClient *classic.AutoscaleSettingsClient

	// alerts management
	ActionRulesClient             *alertsmanagement.ActionRulesClient
	SmartDetectorAlertRulesClient *alertsmanagement.SmartDetectorAlertRulesClient

	// Monitor
	ActionGroupsClient               *classic.ActionGroupsClient
	ActivityLogAlertsClient          *insights.ActivityLogAlertsClient
	AlertRulesClient                 *classic.AlertRulesClient
	DiagnosticSettingsClient         *classic.DiagnosticSettingsClient
	DiagnosticSettingsCategoryClient *classic.DiagnosticSettingsCategoryClient
	LogProfilesClient                *classic.LogProfilesClient
	MetricAlertsClient               *classic.MetricAlertsClient
	ScheduledQueryRulesClient        *classic.ScheduledQueryRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	AADDiagnosticSettingsClient := aad.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AADDiagnosticSettingsClient.Client, o.ResourceManagerAuthorizer)

	AutoscaleSettingsClient := classic.NewAutoscaleSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AutoscaleSettingsClient.Client, o.ResourceManagerAuthorizer)

	ActionRulesClient := alertsmanagement.NewActionRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActionRulesClient.Client, o.ResourceManagerAuthorizer)

	SmartDetectorAlertRulesClient := alertsmanagement.NewSmartDetectorAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SmartDetectorAlertRulesClient.Client, o.ResourceManagerAuthorizer)

	ActionGroupsClient := classic.NewActionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActionGroupsClient.Client, o.ResourceManagerAuthorizer)

	ActivityLogAlertsClient := insights.NewActivityLogAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActivityLogAlertsClient.Client, o.ResourceManagerAuthorizer)

	AlertRulesClient := classic.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AlertRulesClient.Client, o.ResourceManagerAuthorizer)

	DiagnosticSettingsClient := classic.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DiagnosticSettingsClient.Client, o.ResourceManagerAuthorizer)

	DiagnosticSettingsCategoryClient := classic.NewDiagnosticSettingsCategoryClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&DiagnosticSettingsCategoryClient.Client, o.ResourceManagerAuthorizer)

	LogProfilesClient := classic.NewLogProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LogProfilesClient.Client, o.ResourceManagerAuthorizer)

	MetricAlertsClient := classic.NewMetricAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MetricAlertsClient.Client, o.ResourceManagerAuthorizer)

	ScheduledQueryRulesClient := classic.NewScheduledQueryRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ScheduledQueryRulesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AADDiagnosticSettingsClient:      &AADDiagnosticSettingsClient,
		AutoscaleSettingsClient:          &AutoscaleSettingsClient,
		ActionRulesClient:                &ActionRulesClient,
		SmartDetectorAlertRulesClient:    &SmartDetectorAlertRulesClient,
		ActionGroupsClient:               &ActionGroupsClient,
		ActivityLogAlertsClient:          &ActivityLogAlertsClient,
		AlertRulesClient:                 &AlertRulesClient,
		DiagnosticSettingsClient:         &DiagnosticSettingsClient,
		DiagnosticSettingsCategoryClient: &DiagnosticSettingsCategoryClient,
		LogProfilesClient:                &LogProfilesClient,
		MetricAlertsClient:               &MetricAlertsClient,
		ScheduledQueryRulesClient:        &ScheduledQueryRulesClient,
	}
}
