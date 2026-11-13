// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_aad_diagnostic_setting

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAADDiagnosticSettingCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diagnosticsettings.NewDiagnosticSettingID(d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_aad_diagnostic_setting", id.ID())
		}
	}

	// If there is no `enabled` log entry, the PUT will succeed while the next GET will return a 404.
	// Therefore, ensure users has at least one enabled log entry.
	valid := false
	var logs []diagnosticsettings.LogSettings
	if enabledLogs, ok := d.GetOk("enabled_log"); ok && len(enabledLogs.(*pluginsdk.Set).List()) > 0 {
		logs = expandMonitorAADDiagnosticsSettingsEnabledLogs(enabledLogs.(*pluginsdk.Set).List())
		valid = true
	}

	if !valid {
		return fmt.Errorf("at least one of the `log` of the %s should be enabled", id)
	}

	payload := diagnosticsettings.DiagnosticSettingsResource{
		Properties: &diagnosticsettings.DiagnosticSettings{
			Logs: &logs,
		},
	}

	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		payload.Properties.EventHubAuthorizationRuleId = pointer.To(eventHubAuthorizationRuleId)
		payload.Properties.EventHubName = pointer.To(eventHubName)
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		payload.Properties.WorkspaceId = pointer.To(workspaceId)
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		payload.Properties.StorageAccountId = pointer.To(storageAccountId)
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// custom poller is required until https://github.com/Azure/azure-rest-api-specs/issues/38858 is resolved
	pollerType := NewAadDiagnosticSettingCreatePoller(client, id)
	poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	return resourceMonitorAADDiagnosticSettingRead(d, meta)
}

func expandMonitorAADDiagnosticsSettingsEnabledLogs(input []any) []diagnosticsettings.LogSettings {
	results := make([]diagnosticsettings.LogSettings, 0)

	for _, raw := range input {
		if raw == nil {
			continue
		}
		v := raw.(map[string]any)

		logSettings := diagnosticsettings.LogSettings{
			Category: pointer.ToEnum[diagnosticsettings.Category](v["category"].(string)),
			Enabled:  true,
		}

		results = append(results, logSettings)
	}

	return results
}
