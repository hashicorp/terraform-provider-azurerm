// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_diagnostic_setting

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

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
