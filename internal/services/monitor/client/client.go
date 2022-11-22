package client

import (
	"github.com/Azure/azure-sdk-for-go/services/aad/mgmt/2017-04-01/aad"
	"github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2020-10-01/insights"
	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement"
	classic "github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	newActionGroupClient "github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-09-01-preview/insights"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-04-01/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-04-01/datacollectionruleassociations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-04-01/datacollectionrules"
	diagnosticSettingClient "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	diagnosticCategoryClient "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettingscategories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-08-01/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	// AAD
	AADDiagnosticSettingsClient *aad.DiagnosticSettingsClient

	// Autoscale Settings
	AutoscaleSettingsClient *classic.AutoscaleSettingsClient

	// alerts management
	ActionRulesClient             *alertsmanagement.ActionRulesClient
	AlertProcessingRulesClient    *alertprocessingrules.AlertProcessingRulesClient
	SmartDetectorAlertRulesClient *alertsmanagement.SmartDetectorAlertRulesClient

	// Monitor
	ActionGroupsClient                   *newActionGroupClient.ActionGroupsClient
	ActivityLogsClient                   *classic.ActivityLogsClient
	ActivityLogAlertsClient              *insights.ActivityLogAlertsClient
	AlertRulesClient                     *classic.AlertRulesClient
	DataCollectionEndpointsClient        *datacollectionendpoints.DataCollectionEndpointsClient
	DataCollectionRuleAssociationsClient *datacollectionruleassociations.DataCollectionRuleAssociationsClient
	DataCollectionRulesClient            *datacollectionrules.DataCollectionRulesClient
	DiagnosticSettingsClient             *diagnosticSettingClient.DiagnosticSettingsClient
	DiagnosticSettingsCategoryClient     *diagnosticCategoryClient.DiagnosticSettingsCategoriesClient
	LogProfilesClient                    *classic.LogProfilesClient
	MetricAlertsClient                   *classic.MetricAlertsClient
	PrivateLinkScopesClient              *classic.PrivateLinkScopesClient
	PrivateLinkScopedResourcesClient     *classic.PrivateLinkScopedResourcesClient
	ScheduledQueryRulesClient            *classic.ScheduledQueryRulesClient
	ScheduledQueryRulesV2Client          *scheduledqueryrules.ScheduledQueryRulesClient
}

func NewClient(o *common.ClientOptions) *Client {
	AADDiagnosticSettingsClient := aad.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AADDiagnosticSettingsClient.Client, o.ResourceManagerAuthorizer)

	AutoscaleSettingsClient := classic.NewAutoscaleSettingsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AutoscaleSettingsClient.Client, o.ResourceManagerAuthorizer)

	ActionRulesClient := alertsmanagement.NewActionRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActionRulesClient.Client, o.ResourceManagerAuthorizer)

	AlertProcessingRulesClient := alertprocessingrules.NewAlertProcessingRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&AlertProcessingRulesClient.Client, o.ResourceManagerAuthorizer)

	SmartDetectorAlertRulesClient := alertsmanagement.NewSmartDetectorAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SmartDetectorAlertRulesClient.Client, o.ResourceManagerAuthorizer)

	ActionGroupsClient := newActionGroupClient.NewActionGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActionGroupsClient.Client, o.ResourceManagerAuthorizer)

	activityLogsClient := classic.NewActivityLogsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&activityLogsClient.Client, o.ResourceManagerAuthorizer)

	ActivityLogAlertsClient := insights.NewActivityLogAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ActivityLogAlertsClient.Client, o.ResourceManagerAuthorizer)

	AlertRulesClient := classic.NewAlertRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&AlertRulesClient.Client, o.ResourceManagerAuthorizer)

	DataCollectionEndpointsClient := datacollectionendpoints.NewDataCollectionEndpointsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataCollectionEndpointsClient.Client, o.ResourceManagerAuthorizer)

	DataCollectionRuleAssociationsClient := datacollectionruleassociations.NewDataCollectionRuleAssociationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataCollectionRuleAssociationsClient.Client, o.ResourceManagerAuthorizer)

	DataCollectionRulesClient := datacollectionrules.NewDataCollectionRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DataCollectionRulesClient.Client, o.ResourceManagerAuthorizer)

	DiagnosticSettingsClient := diagnosticSettingClient.NewDiagnosticSettingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DiagnosticSettingsClient.Client, o.ResourceManagerAuthorizer)

	DiagnosticSettingsCategoryClient := diagnosticCategoryClient.NewDiagnosticSettingsCategoriesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&DiagnosticSettingsCategoryClient.Client, o.ResourceManagerAuthorizer)

	LogProfilesClient := classic.NewLogProfilesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&LogProfilesClient.Client, o.ResourceManagerAuthorizer)

	MetricAlertsClient := classic.NewMetricAlertsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&MetricAlertsClient.Client, o.ResourceManagerAuthorizer)

	PrivateLinkScopesClient := classic.NewPrivateLinkScopesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PrivateLinkScopesClient.Client, o.ResourceManagerAuthorizer)

	PrivateLinkScopedResourcesClient := classic.NewPrivateLinkScopedResourcesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&PrivateLinkScopedResourcesClient.Client, o.ResourceManagerAuthorizer)

	ScheduledQueryRulesClient := classic.NewScheduledQueryRulesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&ScheduledQueryRulesClient.Client, o.ResourceManagerAuthorizer)

	ScheduledQueryRulesV2Client := scheduledqueryrules.NewScheduledQueryRulesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&ScheduledQueryRulesV2Client.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AADDiagnosticSettingsClient:          &AADDiagnosticSettingsClient,
		AutoscaleSettingsClient:              &AutoscaleSettingsClient,
		ActionRulesClient:                    &ActionRulesClient,
		SmartDetectorAlertRulesClient:        &SmartDetectorAlertRulesClient,
		ActionGroupsClient:                   &ActionGroupsClient,
		ActivityLogsClient:                   &activityLogsClient,
		ActivityLogAlertsClient:              &ActivityLogAlertsClient,
		AlertRulesClient:                     &AlertRulesClient,
		AlertProcessingRulesClient:           &AlertProcessingRulesClient,
		DataCollectionEndpointsClient:        &DataCollectionEndpointsClient,
		DataCollectionRuleAssociationsClient: &DataCollectionRuleAssociationsClient,
		DataCollectionRulesClient:            &DataCollectionRulesClient,
		DiagnosticSettingsClient:             &DiagnosticSettingsClient,
		DiagnosticSettingsCategoryClient:     &DiagnosticSettingsCategoryClient,
		LogProfilesClient:                    &LogProfilesClient,
		MetricAlertsClient:                   &MetricAlertsClient,
		PrivateLinkScopesClient:              &PrivateLinkScopesClient,
		PrivateLinkScopedResourcesClient:     &PrivateLinkScopedResourcesClient,
		ScheduledQueryRulesClient:            &ScheduledQueryRulesClient,
		ScheduledQueryRulesV2Client:          &ScheduledQueryRulesV2Client,
	}
}
