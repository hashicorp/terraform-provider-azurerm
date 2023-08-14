// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azureactivedirectory/2017-04-01/diagnosticsettings"
	authRuleParse "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorAADDiagnosticSetting() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMonitorAADDiagnosticSettingCreate,
		Read:   resourceMonitorAADDiagnosticSettingRead,
		Update: resourceMonitorAADDiagnosticSettingUpdate,
		Delete: resourceMonitorAADDiagnosticSettingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := diagnosticsettings.ParseDiagnosticSettingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MonitorDiagnosticSettingName,
			},

			// When absent, will use the default eventhub, whilst the Diagnostic Setting API will return this property as an empty string. Therefore, it is useless to make this property as Computed.
			"eventhub_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
					"The event hub name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 50 characters, and it must begin and end with a letter or number.",
				),
			},

			"eventhub_authorization_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: authRuleParse.ValidateAuthorizationRuleID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
			},

			"enabled_log": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"retention_policy": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},

									"days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Default:      0,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["enabled_log"].ExactlyOneOf = []string{"enabled_log", "log"}
		resource.Schema["log"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Computed:     true,
			Deprecated:   "`log` has been superseded by `enabled_log` and will be removed in version 4.0 of the AzureRM Provider.",
			ExactlyOneOf: []string{"enabled_log", "log"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
									Default:      0,
								},
							},
						},
					},
				},
			},
		}
	}

	return resource
}

func resourceMonitorAADDiagnosticSettingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := diagnosticsettings.NewDiagnosticSettingID(d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_monitor_aad_diagnostic_setting", id.ID())
	}

	// If there is no `enabled` log entry, the PUT will succeed while the next GET will return a 404.
	// Therefore, ensure users has at least one enabled log entry.
	valid := false
	var logs []diagnosticsettings.LogSettings

	if !features.FourPointOhBeta() {
		if logsRaw, ok := d.GetOk("log"); ok && len(logsRaw.(*pluginsdk.Set).List()) > 0 {
			logs = expandMonitorAADDiagnosticsSettingsLogs(d.Get("log").(*pluginsdk.Set).List())

			for _, v := range logs {
				if v.Enabled {
					valid = true
					break
				}
			}
		}
	}

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

	return resourceMonitorAADDiagnosticSettingRead(d, meta)
}

func resourceMonitorAADDiagnosticSettingUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	if !features.FourPointOhBeta() {
		if d.HasChange("log") {
			logsChanged = true
			logs = expandMonitorAADDiagnosticsSettingsLogs(d.Get("log").(*pluginsdk.Set).List())
			for _, v := range logs {
				if v.Enabled {
					valid = true
					break
				}
			}
		}
	}

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

func resourceMonitorAADDiagnosticSettingRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
				parsedId, err := authRuleParse.ParseAuthorizationRuleIDInsensitively(*props.EventHubAuthorizationRuleId)
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

			if !features.FourPointOhBeta() {
				if err := d.Set("log", flattenMonitorAADDiagnosticLogs(props.Logs)); err != nil {
					return fmt.Errorf("setting `log`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceMonitorAADDiagnosticSettingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := diagnosticsettings.ParseDiagnosticSettingID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	waitForAADDiagnosticSettingToBeGone := waitForAADDiagnosticSettingToBeGonePoller{
		client: client,
		id:     *id,
	}
	initialDelayDuration := 15 * time.Second
	poller := pollers.NewPoller(waitForAADDiagnosticSettingToBeGone, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
	}

	return nil
}

func expandMonitorAADDiagnosticsSettingsLogs(input []interface{}) []diagnosticsettings.LogSettings {
	results := make([]diagnosticsettings.LogSettings, 0)

	for _, raw := range input {
		if raw == nil {
			continue
		}
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		enabled := v["enabled"].(bool)

		policyRaw := v["retention_policy"].([]interface{})[0].(map[string]interface{})
		if len(v["retention_policy"].([]interface{})) == 0 || v["retention_policy"].([]interface{})[0] == nil {
			continue
		}
		retentionDays := policyRaw["days"].(int)
		retentionEnabled := policyRaw["enabled"].(bool)

		results = append(results, diagnosticsettings.LogSettings{
			Category: pointer.To(diagnosticsettings.Category(category)),
			Enabled:  enabled,
			RetentionPolicy: &diagnosticsettings.RetentionPolicy{
				Days:    int64(retentionDays),
				Enabled: retentionEnabled,
			},
		})
	}

	return results
}

func expandMonitorAADDiagnosticsSettingsEnabledLogs(input []interface{}) []diagnosticsettings.LogSettings {
	results := make([]diagnosticsettings.LogSettings, 0)

	for _, raw := range input {
		if raw == nil {
			continue
		}
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		if len(v["retention_policy"].([]interface{})) == 0 || v["retention_policy"].([]interface{})[0] == nil {
			continue
		}

		policyRaw := v["retention_policy"].([]interface{})[0].(map[string]interface{})
		retentionDays := policyRaw["days"].(int)
		retentionEnabled := policyRaw["enabled"].(bool)
		results = append(results, diagnosticsettings.LogSettings{
			Category: pointer.To(diagnosticsettings.Category(category)),
			Enabled:  true,
			RetentionPolicy: &diagnosticsettings.RetentionPolicy{
				Days:    int64(retentionDays),
				Enabled: retentionEnabled,
			},
		})
	}

	return results
}

func flattenMonitorAADDiagnosticLogs(input *[]diagnosticsettings.LogSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		policies := make([]interface{}, 0)
		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			policies = append(policies, map[string]interface{}{
				"days":    int(inputPolicy.Days),
				"enabled": inputPolicy.Enabled,
			})
		}

		category := ""
		if v.Category != nil {
			category = string(*v.Category)
		}
		results = append(results, map[string]interface{}{
			"category":         category,
			"enabled":          v.Enabled,
			"retention_policy": policies,
		})
	}

	return results
}

func flattenMonitorAADDiagnosticEnabledLogs(input *[]diagnosticsettings.LogSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		if !v.Enabled {
			continue
		}

		policies := make([]interface{}, 0)
		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			policies = append(policies, map[string]interface{}{
				"days":    int(inputPolicy.Days),
				"enabled": inputPolicy.Enabled,
			})
		}

		category := ""
		if v.Category != nil {
			category = string(*v.Category)
		}

		results = append(results, map[string]interface{}{
			"category":         category,
			"retention_policy": policies,
		})
	}

	return results
}

var _ pollers.PollerType = waitForAADDiagnosticSettingToBeGonePoller{}

type waitForAADDiagnosticSettingToBeGonePoller struct {
	client *diagnosticsettings.DiagnosticSettingsClient
	id     diagnosticsettings.DiagnosticSettingId
}

func (p waitForAADDiagnosticSettingToBeGonePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("retrieving the deleted %s to check the deletion status: %+v", p.id, err)
		}

		return &pollers.PollResult{
			HttpResponse: &client.Response{
				Response: resp.HttpResponse,
			},
			PollInterval: 15 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	}

	return &pollers.PollResult{
		HttpResponse: &client.Response{
			Response: resp.HttpResponse,
		},
		PollInterval: 15 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
