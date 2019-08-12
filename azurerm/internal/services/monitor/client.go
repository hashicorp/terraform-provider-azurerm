package monitor

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	// Autoscale Settings
	AutoscaleSettingsClient insights.AutoscaleSettingsClient

	// Monitor
	ActionGroupsClient               insights.ActionGroupsClient
	ActivityLogAlertsClient          insights.ActivityLogAlertsClient
	AlertRulesClient                 insights.AlertRulesClient
	DiagnosticSettingsClient         insights.DiagnosticSettingsClient
	DiagnosticSettingsCategoryClient insights.DiagnosticSettingsCategoryClient
	LogProfilesClient                insights.LogProfilesClient
	MetricAlertsClient               insights.MetricAlertsClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.AutoscaleSettingsClient = insights.NewAutoscaleSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AutoscaleSettingsClient.Client, o.ResourceManagerAuthorizer)

	c.ActionGroupsClient = insights.NewActionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ActionGroupsClient.Client, o.ResourceManagerAuthorizer)

	c.ActivityLogAlertsClient = insights.NewActivityLogAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ActivityLogAlertsClient.Client, o.ResourceManagerAuthorizer)

	c.AlertRulesClient = insights.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.AlertRulesClient.Client, o.ResourceManagerAuthorizer)

	c.DiagnosticSettingsClient = insights.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DiagnosticSettingsClient.Client, o.ResourceManagerAuthorizer)

	c.DiagnosticSettingsCategoryClient = insights.NewDiagnosticSettingsCategoryClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.DiagnosticSettingsCategoryClient.Client, o.ResourceManagerAuthorizer)

	c.LogProfilesClient = insights.NewLogProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.LogProfilesClient.Client, o.ResourceManagerAuthorizer)

	c.MetricAlertsClient = insights.NewMetricAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.MetricAlertsClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
