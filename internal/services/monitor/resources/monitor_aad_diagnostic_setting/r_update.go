// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_aad_diagnostic_setting

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAADDiagnosticSettingUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diagnosticsettings.ParseDiagnosticSettingID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	var logs []diagnosticsettings.LogSettings
	logsChanged := false
	valid := false

	if d.HasChange("enabled_log") {
		logsChanged = true
		logs = append(logs, expandMonitorAADDiagnosticsSettingsEnabledLogs(d.Get("enabled_log").(*pluginsdk.Set).List())...)
		valid = true
	}

	if !logsChanged && existing.Model.Properties != nil && existing.Model.Properties.Logs != nil {
		logs = *existing.Model.Properties.Logs
		for _, v := range logs {
			if v.Enabled {
				valid = true
				break
			}
		}
	}

	if !valid {
		return fmt.Errorf("at least one of the `log` of the %s should be enabled", id)
	}

	properties := diagnosticsettings.DiagnosticSettingsResource{
		Properties: &diagnosticsettings.DiagnosticSettings{
			Logs: &logs,
		},
	}

	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		properties.Properties.EventHubAuthorizationRuleId = pointer.To(eventHubAuthorizationRuleId)
		properties.Properties.EventHubName = pointer.To(eventHubName)
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		properties.Properties.WorkspaceId = pointer.To(workspaceId)
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		properties.Properties.StorageAccountId = pointer.To(storageAccountId)
	}

	if _, err := client.CreateOrUpdate(ctx, *id, properties); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return resourceMonitorAADDiagnosticSettingRead(d, meta)
}
