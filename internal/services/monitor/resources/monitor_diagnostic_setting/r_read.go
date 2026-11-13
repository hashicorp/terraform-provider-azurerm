// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_diagnostic_setting

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

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

	d.Set("name", id.DiagnosticSettingName)
	resourceUri := id.ResourceUri
	if v, err := commonids.ParseKustoClusterIDInsensitively(resourceUri); err == nil {
		resourceUri = v.ID()
	}
	d.Set("target_resource_id", resourceUri)

	if model := resp.Model; model != nil {
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
			if resp.Model.Properties.LogAnalyticsDestinationType != nil && *resp.Model.Properties.LogAnalyticsDestinationType != "" {
				logAnalyticsDestinationType = *resp.Model.Properties.LogAnalyticsDestinationType
			}
			d.Set("log_analytics_destination_type", logAnalyticsDestinationType)

			if err = d.Set("enabled_log", flattenMonitorDiagnosticEnabledLogs(resp.Model.Properties.Logs)); err != nil {
				return fmt.Errorf("setting `enabled_log`: %+v", err)
			}

			if err = d.Set("enabled_metric", flattenMonitorDiagnosticEnabledMetrics(resp.Model.Properties.Metrics)); err != nil {
				return fmt.Errorf("setting `enabled_metric`: %+v", err)
			}
		}
	}

	return nil
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

func ParseMonitorDiagnosticId(monitorId string) (*diagnosticsettings.ScopedDiagnosticSettingId, error) {
	v := strings.Split(monitorId, "|")
	if len(v) != 2 {
		return nil, fmt.Errorf("expected the Monitor Diagnostics ID to be in the format `{resourceId}|{name}` but got %d segments", len(v))
	}

	// TODO: this can become a Composite Resource ID once https://github.com/hashicorp/go-azure-helpers/pull/208 is released
	identifier := diagnosticsettings.NewScopedDiagnosticSettingID(v[0], v[1])
	return &identifier, nil
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
