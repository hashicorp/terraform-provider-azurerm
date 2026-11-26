// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_aad_diagnostic_setting

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAADDiagnosticSettingRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diagnosticsettings.ParseDiagnosticSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.DiagnosticSettingName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("eventhub_name", props.EventHubName)

			eventhubAuthorizationRuleId := ""
			if props.EventHubAuthorizationRuleId != nil && *props.EventHubAuthorizationRuleId != "" {
				parsedId, err := authorizationrulesnamespaces.ParseAuthorizationRuleIDInsensitively(*props.EventHubAuthorizationRuleId)
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

			storageAccountId := ""
			if props.StorageAccountId != nil && *props.StorageAccountId != "" {
				parsedId, err := commonids.ParseStorageAccountIDInsensitively(*props.StorageAccountId)
				if err != nil {
					return err
				}

				storageAccountId = parsedId.ID()
			}
			d.Set("storage_account_id", storageAccountId)

			if err := d.Set("enabled_log", flattenMonitorAADDiagnosticEnabledLogs(props.Logs)); err != nil {
				return fmt.Errorf("setting `enabled_log`: %+v", err)
			}
		}
	}

	return nil
}

func flattenMonitorAADDiagnosticEnabledLogs(input *[]diagnosticsettings.LogSettings) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		if !v.Enabled {
			continue
		}

		category := ""
		if v.Category != nil {
			category = string(*v.Category)
		}

		result := map[string]any{
			"category": category,
		}

		results = append(results, result)
	}

	return results
}

var _ pollers.PollerType = &waitForAADDiagnosticSettingToBeGonePoller{}

type waitForAADDiagnosticSettingToBeGonePoller struct {
	client                     *diagnosticsettings.DiagnosticSettingsClient
	id                         diagnosticsettings.DiagnosticSettingId
	continuousTargetOccurrence int
}

func (p *waitForAADDiagnosticSettingToBeGonePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil && !response.WasNotFound(resp.HttpResponse) {
		return nil, fmt.Errorf("retrieving the deleted %s to check the deletion status: %+v", p.id, err)
	}

	if response.WasNotFound(resp.HttpResponse) {
		if p.continuousTargetOccurrence >= 3 {
			return &pollers.PollResult{
				HttpResponse: &client.Response{
					Response: resp.HttpResponse,
				},
				PollInterval: 15 * time.Second,
				Status:       pollers.PollingStatusSucceeded,
			}, nil
		}
		p.continuousTargetOccurrence++
	} else {
		p.continuousTargetOccurrence = 0
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}

var _ pollers.PollerType = &aadDiagnosticSettingCreatePoller{}

type aadDiagnosticSettingCreatePoller struct {
	client *diagnosticsettings.DiagnosticSettingsClient
	id     diagnosticsettings.DiagnosticSettingId
}

func NewAadDiagnosticSettingCreatePoller(client *diagnosticsettings.DiagnosticSettingsClient, id diagnosticsettings.DiagnosticSettingId) *aadDiagnosticSettingCreatePoller {
	return &aadDiagnosticSettingCreatePoller{
		client: client,
		id:     id,
	}
}

func (p *aadDiagnosticSettingCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &pollers.PollResult{
				HttpResponse: &client.Response{
					Response: resp.HttpResponse,
				},
				PollInterval: 15 * time.Second,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}, nil
}
