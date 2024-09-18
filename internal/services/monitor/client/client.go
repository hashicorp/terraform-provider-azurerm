// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2015-04-01/activitylogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	scheduledqueryrules2018 "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2019-10-17-preview/privatelinkscopedresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2020-10-01/activitylogalertsapis"
	diagnosticSettingClient "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	diagnosticCategoryClient "github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettingscategories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionruleassociations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	// AAD
	AADDiagnosticSettingsClient *diagnosticsettings.DiagnosticSettingsClient

	// Autoscale Settings
	AutoscaleSettingsClient *autoscalesettings.AutoScaleSettingsClient

	// alerts management
	AlertProcessingRulesClient    *alertprocessingrules.AlertProcessingRulesClient
	SmartDetectorAlertRulesClient *smartdetectoralertrules.SmartDetectorAlertRulesClient

	// Monitor
	ActionGroupsClient                   *actiongroupsapis.ActionGroupsAPIsClient
	ActivityLogsClient                   *activitylogs.ActivityLogsClient
	ActivityLogAlertsClient              *activitylogalertsapis.ActivityLogAlertsAPIsClient
	AlertPrometheusRuleGroupClient       *prometheusrulegroups.PrometheusRuleGroupsClient
	DataCollectionEndpointsClient        *datacollectionendpoints.DataCollectionEndpointsClient
	DataCollectionRuleAssociationsClient *datacollectionruleassociations.DataCollectionRuleAssociationsClient
	DataCollectionRulesClient            *datacollectionrules.DataCollectionRulesClient
	DiagnosticSettingsClient             *diagnosticSettingClient.DiagnosticSettingsClient
	DiagnosticSettingsCategoryClient     *diagnosticCategoryClient.DiagnosticSettingsCategoriesClient
	MetricAlertsClient                   *metricalerts.MetricAlertsClient
	PrivateLinkScopesClient              *privatelinkscopesapis.PrivateLinkScopesAPIsClient
	PrivateLinkScopedResourcesClient     *privatelinkscopedresources.PrivateLinkScopedResourcesClient
	ScheduledQueryRulesClient            *scheduledqueryrules2018.ScheduledQueryRulesClient
	ScheduledQueryRulesV2Client          *scheduledqueryrules.ScheduledQueryRulesClient
	WorkspacesClient                     *azuremonitorworkspaces.AzureMonitorWorkspacesClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	aadDiagnosticSettingsClient, err := diagnosticsettings.NewDiagnosticSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AAD DiagnosticsSettings client: %+v", err)
	}
	o.Configure(aadDiagnosticSettingsClient.Client, o.Authorizers.ResourceManager)

	AutoscaleSettingsClient, err := autoscalesettings.NewAutoScaleSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Autoscale Settings client: %+v", err)
	}
	o.Configure(AutoscaleSettingsClient.Client, o.Authorizers.ResourceManager)

	alertProcessingRulesClient, err := alertprocessingrules.NewAlertProcessingRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building AlertProcessingRules client: %+v", err)
	}
	o.Configure(alertProcessingRulesClient.Client, o.Authorizers.ResourceManager)

	SmartDetectorAlertRulesClient, err := smartdetectoralertrules.NewSmartDetectorAlertRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Smart Detector Alert Rules client: %+v", err)
	}
	o.Configure(SmartDetectorAlertRulesClient.Client, o.Authorizers.ResourceManager)

	ActionGroupsClient, err := actiongroupsapis.NewActionGroupsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Action Groups client: %+v", err)
	}
	o.Configure(ActionGroupsClient.Client, o.Authorizers.ResourceManager)

	activityLogsClient, err := activitylogs.NewActivityLogsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Activity Logs client: %+v", err)
	}
	o.Configure(activityLogsClient.Client, o.Authorizers.ResourceManager)

	ActivityLogAlertsClient, err := activitylogalertsapis.NewActivityLogAlertsAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Activity Log Alerts client: %+v", err)
	}
	o.Configure(ActivityLogAlertsClient.Client, o.Authorizers.ResourceManager)

	alertPrometheusRuleGroupClient, err := prometheusrulegroups.NewPrometheusRuleGroupsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building PrometheusRuleGroups client: %+v", err)
	}
	o.Configure(alertPrometheusRuleGroupClient.Client, o.Authorizers.ResourceManager)

	DataCollectionEndpointsClient, err := datacollectionendpoints.NewDataCollectionEndpointsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Collection Endpoints client: %+v", err)
	}
	o.Configure(DataCollectionEndpointsClient.Client, o.Authorizers.ResourceManager)

	DataCollectionRuleAssociationsClient, err := datacollectionruleassociations.NewDataCollectionRuleAssociationsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Collection Rule Associations client: %+v", err)
	}
	o.Configure(DataCollectionRuleAssociationsClient.Client, o.Authorizers.ResourceManager)

	DataCollectionRulesClient, err := datacollectionrules.NewDataCollectionRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Data Collection Rules client: %+v", err)
	}
	o.Configure(DataCollectionRulesClient.Client, o.Authorizers.ResourceManager)

	DiagnosticSettingsClient, err := diagnosticSettingClient.NewDiagnosticSettingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Diagnostic Settings client: %+v", err)
	}
	o.Configure(DiagnosticSettingsClient.Client, o.Authorizers.ResourceManager)

	DiagnosticSettingsCategoryClient, err := diagnosticCategoryClient.NewDiagnosticSettingsCategoriesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Diagnostic Settings Categories client: %+v", err)
	}
	o.Configure(DiagnosticSettingsCategoryClient.Client, o.Authorizers.ResourceManager)

	MetricAlertsClient, err := metricalerts.NewMetricAlertsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Metric Alerts client: %+v", err)
	}
	o.Configure(MetricAlertsClient.Client, o.Authorizers.ResourceManager)

	PrivateLinkScopesClient, err := privatelinkscopesapis.NewPrivateLinkScopesAPIsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Private Link Scopes client: %+v", err)
	}
	o.Configure(PrivateLinkScopesClient.Client, o.Authorizers.ResourceManager)

	PrivateLinkScopedResourcesClient, err := privatelinkscopedresources.NewPrivateLinkScopedResourcesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Private Link Scoped Resources client: %+v", err)
	}
	o.Configure(PrivateLinkScopedResourcesClient.Client, o.Authorizers.ResourceManager)

	ScheduledQueryRulesClient, err := scheduledqueryrules2018.NewScheduledQueryRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Scheduled Query Rules client: %+v", err)
	}
	o.Configure(ScheduledQueryRulesClient.Client, o.Authorizers.ResourceManager)

	ScheduledQueryRulesV2Client, err := scheduledqueryrules.NewScheduledQueryRulesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Scheduled Query Rules V2 client: %+v", err)
	}
	o.Configure(ScheduledQueryRulesV2Client.Client, o.Authorizers.ResourceManager)

	WorkspacesClient, err := azuremonitorworkspaces.NewAzureMonitorWorkspacesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Workspaces client: %+v", err)
	}
	o.Configure(WorkspacesClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		AADDiagnosticSettingsClient:          aadDiagnosticSettingsClient,
		AutoscaleSettingsClient:              AutoscaleSettingsClient,
		SmartDetectorAlertRulesClient:        SmartDetectorAlertRulesClient,
		ActionGroupsClient:                   ActionGroupsClient,
		ActivityLogsClient:                   activityLogsClient,
		ActivityLogAlertsClient:              ActivityLogAlertsClient,
		AlertPrometheusRuleGroupClient:       alertPrometheusRuleGroupClient,
		AlertProcessingRulesClient:           alertProcessingRulesClient,
		DataCollectionEndpointsClient:        DataCollectionEndpointsClient,
		DataCollectionRuleAssociationsClient: DataCollectionRuleAssociationsClient,
		DataCollectionRulesClient:            DataCollectionRulesClient,
		DiagnosticSettingsClient:             DiagnosticSettingsClient,
		DiagnosticSettingsCategoryClient:     DiagnosticSettingsCategoryClient,
		MetricAlertsClient:                   MetricAlertsClient,
		PrivateLinkScopesClient:              PrivateLinkScopesClient,
		PrivateLinkScopedResourcesClient:     PrivateLinkScopedResourcesClient,
		ScheduledQueryRulesClient:            ScheduledQueryRulesClient,
		ScheduledQueryRulesV2Client:          ScheduledQueryRulesV2Client,
		WorkspacesClient:                     WorkspacesClient,
	}, nil
}
