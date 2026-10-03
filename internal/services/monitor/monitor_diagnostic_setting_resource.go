// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor

//go:generate go run ../../tools/generator-tests resourceidentity -properties "resource_uri:target_resource_id,name"

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorDiagnosticSetting() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorDiagnosticSettingCreate,
		Read:   resourceMonitorDiagnosticSettingRead,
		Update: resourceMonitorDiagnosticSettingUpdate,
		Delete: resourceMonitorDiagnosticSettingDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&diagnosticsettings.ScopedDiagnosticSettingId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&diagnosticsettings.ScopedDiagnosticSettingId{}),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DiagnosticSettingV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MonitorDiagnosticSettingName,
			},

			"target_resource_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					azure.ValidateResourceID,
					commonids.ValidateManagementGroupID,
				),
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"eventhub_authorization_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: authorizationrulesnamespaces.ValidateAuthorizationRuleID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
			},

			"partner_solution_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
			},

			"log_analytics_destination_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true, // azignore:AZS007 - pre-existing violation
				ValidateFunc: validation.StringInSlice([]string{
					"Dedicated",
					"AzureDiagnostics", // Not documented in azure API, but some resource has skew. See: https://github.com/Azure/azure-rest-api-specs/issues/9281
				}, false),
			},

			"enabled_log": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"enabled_log", "enabled_metric"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"category_group": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				Set: resourceMonitorDiagnosticLogSettingHash,
			},

			"enabled_metric": {
				Type:         pluginsdk.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"enabled_log", "enabled_metric"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceMonitorDiagnosticSettingCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diagnosticsettings.NewScopedDiagnosticSettingID(d.Get("target_resource_id").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor Diagnostic Setting %q for Resource %q: %s", id.DiagnosticSettingName, id.ResourceUri, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_diagnostic_setting", id.ID())
		}
	}

	var logs []diagnosticsettings.DiagnosticsLogSettings
	hasEnabledLogs := false
	if enabledLogs, ok := d.GetOk("enabled_log"); ok {
		enabledLogsList := enabledLogs.(*pluginsdk.Set).List()
		if len(enabledLogsList) > 0 {
			expandEnabledLogs, err := expandMonitorDiagnosticsSettingsEnabledLogs(enabledLogsList)
			if err != nil {
				return fmt.Errorf("expanding enabled_log: %+v", err)
			}
			logs = *expandEnabledLogs
			hasEnabledLogs = true
		}
	}

	// if no logs/metrics are enabled the API "creates" but 404's on Read
	var metrics []diagnosticsettings.DiagnosticsMetricSettings
	hasEnabledMetrics := false

	if enabledMetrics, ok := d.GetOk("enabled_metric"); ok {
		enabledMetricsList := enabledMetrics.(*pluginsdk.Set).List()
		if len(enabledMetricsList) > 0 {
			metrics = expandMonitorDiagnosticsSettingsEnabledMetrics(enabledMetricsList)
			hasEnabledMetrics = true
		}
	}

	if !hasEnabledMetrics && !hasEnabledLogs {
		return fmt.Errorf("at least one type of Log or Metric must be enabled")
	}

	parameters := diagnosticsettings.DiagnosticSettingsResource{
		Properties: &diagnosticsettings.DiagnosticSettings{
			Logs:    &logs,
			Metrics: &metrics,
		},
	}

	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		parameters.Properties.EventHubAuthorizationRuleId = pointer.To(eventHubAuthorizationRuleId)
		parameters.Properties.EventHubName = pointer.To(eventHubName)
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		parameters.Properties.WorkspaceId = pointer.To(workspaceId)
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		parameters.Properties.StorageAccountId = pointer.To(storageAccountId)
	}

	partnerSolutionId := d.Get("partner_solution_id").(string)
	if partnerSolutionId != "" {
		parameters.Properties.MarketplacePartnerId = pointer.To(partnerSolutionId)
	}

	if v := d.Get("log_analytics_destination_type").(string); v != "" {
		parameters.Properties.LogAnalyticsDestinationType = &v
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating Monitor Diagnostics Setting %q for Resource %q: %+v", id.DiagnosticSettingName, id.ResourceUri, err)
	}

	// https://github.com/Azure/azure-rest-api-specs/issues/30249
	log.Printf("[DEBUG] Waiting for Monitor Diagnostic Setting %q for Resource %q to become ready", id.DiagnosticSettingName, id.ResourceUri)
	poller := custompollers.NewEventualConsistencyPoller(3, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              5 * time.Second,
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for Monitor Diagnostic Setting %q for Resource %q to become ready: %s", id.DiagnosticSettingName, id.ResourceUri, err)
	}

	d.SetId(id.ID())

	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return fmt.Errorf("setting resource identity: %w", err)
	}

	return resourceMonitorDiagnosticSettingRead(d, meta)
}

func resourceMonitorDiagnosticSettingUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Monitor Diagnostics Setting %q for Resource %q: %+v", id.DiagnosticSettingName, id.ResourceUri, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("unexpected null model of Monitor Diagnostics Setting %q for Resource %q", id.DiagnosticSettingName, id.ResourceUri)
	}

	var logs []diagnosticsettings.DiagnosticsLogSettings
	hasEnabledLogs := false

	if d.HasChange("enabled_log") {
		enabledLogs := d.Get("enabled_log").(*pluginsdk.Set).List()
		log.Printf("[DEBUG] enabled_logs: %+v", enabledLogs)
		if len(enabledLogs) > 0 {
			expandEnabledLogs, err := expandMonitorDiagnosticsSettingsEnabledLogs(enabledLogs)
			if err != nil {
				return fmt.Errorf("expanding enabled_log: %+v", err)
			}
			logs = *expandEnabledLogs
			hasEnabledLogs = true
		} else if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Logs != nil {
			// if the enabled_log is updated to empty, we disable the log explicitly
			for _, v := range *existing.Model.Properties.Logs {
				disabledLog := v
				disabledLog.Enabled = false
				logs = append(logs, disabledLog)
			}
		}
	} else if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Logs != nil {
		logs = *existing.Model.Properties.Logs
		for _, v := range logs {
			if v.Enabled {
				hasEnabledLogs = true
			}
		}
	}

	var metrics []diagnosticsettings.DiagnosticsMetricSettings
	hasEnabledMetrics := false

	if d.HasChange("enabled_metric") {
		enabledMetrics := d.Get("enabled_metric").(*pluginsdk.Set).List()
		if len(enabledMetrics) > 0 {
			metrics = expandMonitorDiagnosticsSettingsEnabledMetrics(enabledMetrics)
			hasEnabledMetrics = true
		} else if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Metrics != nil {
			// if the enabled_metric is updated to empty, we disable the metric explicitly
			for _, v := range *existing.Model.Properties.Metrics {
				disabledMetric := v
				disabledMetric.Enabled = false
				metrics = append(metrics, disabledMetric)
			}
		}
	} else if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.Metrics != nil {
		metrics = *existing.Model.Properties.Metrics
		for _, v := range metrics {
			if v.Enabled {
				hasEnabledMetrics = true
			}
		}
	}

	// if no logs/metrics are enabled the API "creates" but 404's on Read
	if !hasEnabledMetrics && !hasEnabledLogs {
		return fmt.Errorf("at least one type of Log or Metric must be enabled")
	}

	parameters := diagnosticsettings.DiagnosticSettingsResource{
		Properties: &diagnosticsettings.DiagnosticSettings{
			Logs:    &logs,
			Metrics: &metrics,
		},
	}

	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		parameters.Properties.EventHubAuthorizationRuleId = pointer.To(eventHubAuthorizationRuleId)
		parameters.Properties.EventHubName = pointer.To(eventHubName)
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		parameters.Properties.WorkspaceId = pointer.To(workspaceId)
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		parameters.Properties.StorageAccountId = pointer.To(storageAccountId)
	}

	partnerSolutionId := d.Get("partner_solution_id").(string)
	if partnerSolutionId != "" {
		parameters.Properties.MarketplacePartnerId = pointer.To(partnerSolutionId)
	}

	if v := d.Get("log_analytics_destination_type").(string); v != "" {
		parameters.Properties.LogAnalyticsDestinationType = &v
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating Monitor Diagnostics Setting %q for Resource %q: %+v", id.DiagnosticSettingName, id.ResourceUri, err)
	}
	return resourceMonitorDiagnosticSettingRead(d, meta)
}

func resourceMonitorDiagnosticSettingRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] Monitor Diagnostics Setting %q was not found for Resource %q - removing from state!", id.DiagnosticSettingName, id.ResourceUri)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Monitor Diagnostics Setting %q for Resource %q: %+v", id.DiagnosticSettingName, id.ResourceUri, err)
	}

	return resourceMonitorDiagnosticSettingFlatten(d, id, resp.Model)
}

func resourceMonitorDiagnosticSettingFlatten(d *pluginsdk.ResourceData, id *diagnosticsettings.ScopedDiagnosticSettingId, model *diagnosticsettings.DiagnosticSettingsResource) error {
	d.Set("name", id.DiagnosticSettingName)
	resourceUri := id.ResourceUri
	if v, err := commonids.ParseKustoClusterIDInsensitively(resourceUri); err == nil {
		resourceUri = v.ID()
	}
	d.Set("target_resource_id", resourceUri)

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("eventhub_name", props.EventHubName)
			eventhubAuthorizationRuleId := ""
			if props.EventHubAuthorizationRuleId != nil && *props.EventHubAuthorizationRuleId != "" {
				authRuleId := pointer.From(props.EventHubAuthorizationRuleId)
				parsedId, err := authorizationrulesnamespaces.ParseAuthorizationRuleIDInsensitively(authRuleId)
				if err != nil {
					return err
				}
				eventhubAuthorizationRuleId = parsedId.ID()
			}
			d.Set("eventhub_authorization_rule_id", eventhubAuthorizationRuleId)

			workspaceId := ""
			if props.WorkspaceId != nil && *props.WorkspaceId != "" {
				parsedId, err := workspaces.ParseWorkspaceIDInsensitively(*props.WorkspaceId)
				if err != nil {
					return err
				}

				workspaceId = parsedId.ID()
			}
			d.Set("log_analytics_workspace_id", workspaceId)

			if props.StorageAccountId != nil && *props.StorageAccountId != "" {
				parsedId, err := commonids.ParseStorageAccountIDInsensitively(*props.StorageAccountId)
				if err != nil {
					return err
				}

				d.Set("storage_account_id", parsedId.ID())
			}

			if props.MarketplacePartnerId != nil && *props.MarketplacePartnerId != "" {
				d.Set("partner_solution_id", props.MarketplacePartnerId)
			}

			logAnalyticsDestinationType := ""
			if model.Properties.LogAnalyticsDestinationType != nil && *model.Properties.LogAnalyticsDestinationType != "" {
				logAnalyticsDestinationType = *model.Properties.LogAnalyticsDestinationType
			}
			d.Set("log_analytics_destination_type", logAnalyticsDestinationType)

			if err := d.Set("enabled_log", flattenMonitorDiagnosticEnabledLogs(model.Properties.Logs)); err != nil {
				return fmt.Errorf("setting `enabled_log`: %+v", err)
			}

			if err := d.Set("enabled_metric", flattenMonitorDiagnosticEnabledMetrics(model.Properties.Metrics)); err != nil {
				return fmt.Errorf("setting `enabled_metric`: %+v", err)
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceMonitorDiagnosticSettingDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor Diagnostics Setting %q for Resource %q: %+v", id.DiagnosticSettingName, id.ResourceUri, err)
		}
	}

	// API appears to be eventually consistent (identified during tainting this resource)
	log.Printf("[DEBUG] Waiting for Monitor Diagnostic Setting %q for Resource %q to disappear", id.DiagnosticSettingName, id.ResourceUri)
	poller := custompollers.NewEventualConsistencyPoller(5, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, *id)
		return resp.HttpResponse, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:         15 * time.Second,
		TargetStatusCode: pointer.To(http.StatusNotFound),
	})
	if err = poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for Monitor Diagnostic Setting %q for Resource %q to disappear: %s", id.DiagnosticSettingName, id.ResourceUri, err)
	}

	return nil
}

func expandMonitorDiagnosticsSettingsEnabledLogs(input []any) (*[]diagnosticsettings.DiagnosticsLogSettings, error) {
	results := make([]diagnosticsettings.DiagnosticsLogSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		category := v["category"].(string)
		categoryGroup := v["category_group"].(string)

		output := diagnosticsettings.DiagnosticsLogSettings{
			Enabled: true,
		}

		switch {
		case category != "":
			output.Category = pointer.To(category)
		case categoryGroup != "":
			output.CategoryGroup = pointer.To(categoryGroup)
		default:
			return nil, fmt.Errorf("exactly one of `category` or `category_group` must be specified")
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenMonitorDiagnosticEnabledLogs(input *[]diagnosticsettings.DiagnosticsLogSettings) []any {
	enabledLogs := make([]any, 0)
	if input == nil {
		return enabledLogs
	}

	for _, v := range *input {
		output := make(map[string]any)

		if !v.Enabled {
			continue
		}

		output["category"] = pointer.From(v.Category)

		output["category_group"] = pointer.From(v.CategoryGroup)

		enabledLogs = append(enabledLogs, output)
	}
	return enabledLogs
}

func flattenMonitorDiagnosticEnabledMetrics(input *[]diagnosticsettings.DiagnosticsMetricSettings) []any {
	enabledLogs := make([]any, 0)
	if input == nil {
		return enabledLogs
	}

	for _, v := range *input {
		output := make(map[string]any)

		if !v.Enabled {
			continue
		}

		output["category"] = pointer.From(v.Category)
		enabledLogs = append(enabledLogs, output)
	}
	return enabledLogs
}

func expandMonitorDiagnosticsSettingsEnabledMetrics(input []any) []diagnosticsettings.DiagnosticsMetricSettings {
	results := make([]diagnosticsettings.DiagnosticsMetricSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]any)

		output := diagnosticsettings.DiagnosticsMetricSettings{
			Category: pointer.To(v["category"].(string)),
			Enabled:  true,
		}

		results = append(results, output)
	}

	return results
}

func ParseMonitorDiagnosticId(monitorId string) (*diagnosticsettings.ScopedDiagnosticSettingId, error) {
	if strings.Contains(monitorId, "|") {
		v := strings.Split(monitorId, "|")
		if len(v) != 2 {
			return nil, fmt.Errorf("expected the Monitor Diagnostics ID to be in the format `{resourceId}|{name}` but got %d segments", len(v))
		}

		identifier := diagnosticsettings.NewScopedDiagnosticSettingID(v[0], v[1])
		return &identifier, nil
	}

	if !strings.HasPrefix(monitorId, "/") {
		monitorId = "/" + monitorId
	}

	return diagnosticsettings.ParseScopedDiagnosticSettingIDInsensitively(monitorId)
}

func resourceMonitorDiagnosticLogSettingHash(input any) int {
	var buf bytes.Buffer
	if rawData, ok := input.(map[string]any); ok {
		if category, ok := rawData["category"]; ok {
			fmt.Fprintf(&buf, "%s-", category.(string))
		}
		if categoryGroup, ok := rawData["category_group"]; ok {
			fmt.Fprintf(&buf, "%s-", categoryGroup.(string))
		}
	}
	return pluginsdk.HashString(buf.String())
}
