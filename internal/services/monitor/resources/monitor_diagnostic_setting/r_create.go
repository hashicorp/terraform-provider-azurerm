// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_diagnostic_setting

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorDiagnosticSettingCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diagnosticsettings.NewScopedDiagnosticSettingID(d.Get("target_resource_id").(string), d.Get("name").(string))
	resourceId := fmt.Sprintf("%s|%s", id.ResourceUri, id.DiagnosticSettingName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor Diagnostic Setting %q for Resource %q: %s", id.DiagnosticSettingName, id.ResourceUri, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_diagnostic_setting", resourceId)
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

	d.SetId(resourceId)

	return resourceMonitorDiagnosticSettingRead(d, meta)
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
