package monitor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	authRuleParse "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-05-01-preview/diagnosticsettings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	storageParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorDiagnosticSetting() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorDiagnosticSettingCreateUpdate,
		Read:   resourceMonitorDiagnosticSettingRead,
		Update: resourceMonitorDiagnosticSettingCreateUpdate,
		Delete: resourceMonitorDiagnosticSettingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := ParseMonitorDiagnosticId(id)
			return err
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"eventhub_authorization_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: authRuleParse.ValidateAuthorizationRuleID,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"partner_solution_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"log_analytics_destination_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"Dedicated",
					"AzureDiagnostics", // Not documented in azure API, but some resource has skew. See: https://github.com/Azure/azure-rest-api-specs/issues/9281
				}, false),
			},

			"log": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"category_group": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"retention_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},

									"days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
						},
					},
				},
				Set: resourceMonitorDiagnosticLogSettingHash,
			},

			"metric": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
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
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Required: true,
									},

									"days": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
						},
					},
				},
				Set: resourceMonitorDiagnosticMetricsSettingHash,
			},
		},
	}
}

func resourceMonitorDiagnosticSettingCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Diagnostic Settings.")

	name := d.Get("name").(string)
	actualResourceId := d.Get("target_resource_id").(string)
	diagnosticSettingId := diagnosticsettings.NewScopedDiagnosticSettingID(actualResourceId, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, diagnosticSettingId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor Diagnostic Setting %q for Resource %q: %s", diagnosticSettingId.Name, diagnosticSettingId.ResourceUri, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_monitor_diagnostic_setting", *existing.Model.Id)
		}
	}

	logsRaw := d.Get("log").(*pluginsdk.Set).List()
	logs := expandMonitorDiagnosticsSettingsLogs(logsRaw)
	metricsRaw := d.Get("metric").(*pluginsdk.Set).List()
	metrics := expandMonitorDiagnosticsSettingsMetrics(metricsRaw)

	// if no blocks are specified  the API "creates" but 404's on Read
	if len(logs) == 0 && len(metrics) == 0 {
		return fmt.Errorf("At least one `log` or `metric` block must be specified")
	}

	// also if there's none enabled
	valid := false
	for _, v := range logs {
		if v.Enabled {
			valid = true
			break
		}
	}
	if !valid {
		for _, v := range metrics {
			if v.Enabled {
				valid = true
				break
			}
		}
	}

	if !valid {
		return fmt.Errorf("At least one `log` or `metric` must be enabled")
	}

	parameters := diagnosticsettings.DiagnosticSettingsResource{
		Properties: &diagnosticsettings.DiagnosticSettings{
			Logs:    &logs,
			Metrics: &metrics,
		},
	}

	valid = false
	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		parameters.Properties.EventHubAuthorizationRuleId = utils.String(eventHubAuthorizationRuleId)
		parameters.Properties.EventHubName = utils.String(eventHubName)
		valid = true
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		parameters.Properties.WorkspaceId = utils.String(workspaceId)
		valid = true
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		parameters.Properties.StorageAccountId = utils.String(storageAccountId)
		valid = true
	}

	partnerSolutionId := d.Get("partner_solution_id").(string)
	if partnerSolutionId != "" {
		parameters.Properties.MarketplacePartnerId = utils.String(partnerSolutionId)
		valid = true
	}

	if v := d.Get("log_analytics_destination_type").(string); v != "" {
		parameters.Properties.LogAnalyticsDestinationType = &v
	}

	if !valid {
		return fmt.Errorf("either a `eventhub_authorization_rule_id`, `log_analytics_workspace_id`, `partner_solution_id` or `storage_account_id` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, diagnosticSettingId, parameters); err != nil {
		return fmt.Errorf("creating Monitor Diagnostics Setting %q for Resource %q: %+v", name, actualResourceId, err)
	}

	read, err := client.Get(ctx, diagnosticSettingId)
	if err != nil {
		return err
	}
	if read.Model == nil && read.Model.Id == nil {
		return fmt.Errorf("Cannot read ID for Monitor Diagnostics %q for Resource ID %q", diagnosticSettingId.Name, diagnosticSettingId.ResourceUri)
	}

	d.SetId(fmt.Sprintf("%s|%s", actualResourceId, name))

	return resourceMonitorDiagnosticSettingRead(d, meta)
}

func resourceMonitorDiagnosticSettingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	actualResourceId := id.ResourceUri
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] Monitor Diagnostics Setting %q was not found for Resource %q - removing from state!", id.Name, actualResourceId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Monitor Diagnostics Setting %q for Resource %q: %+v", id.Name, actualResourceId, err)
	}

	d.Set("name", id.Name)
	d.Set("target_resource_id", id.ResourceUri)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("eventhub_name", props.EventHubName)
			eventhubAuthorizationRuleId := ""
			if props.EventHubAuthorizationRuleId != nil && *props.EventHubAuthorizationRuleId != "" {
				authRuleId := utils.NormalizeNilableString(props.EventHubAuthorizationRuleId)
				parsedId, err := authRuleParse.ParseAuthorizationRuleIDInsensitively(authRuleId)
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
				parsedId, err := storageParse.StorageAccountID(*props.StorageAccountId)
				if err != nil {
					return err
				}

				storageAccountId = parsedId.ID()
				d.Set("storage_account_id", storageAccountId)
			}

			partnerSolutionId := ""
			if props.MarketplacePartnerId != nil && *props.MarketplacePartnerId != "" {
				partnerSolutionId = *props.MarketplacePartnerId
				d.Set("partner_solution_id", partnerSolutionId)
			}

			d.Set("log_analytics_destination_type", resp.Model.Properties.LogAnalyticsDestinationType)

			if err := d.Set("log", flattenMonitorDiagnosticLogs(resp.Model.Properties.Logs)); err != nil {
				return fmt.Errorf("setting `log`: %+v", err)
			}

			if err := d.Set("metric", flattenMonitorDiagnosticMetrics(resp.Model.Properties.Metrics)); err != nil {
				return fmt.Errorf("setting `metric`: %+v", err)
			}
		}
	}

	return nil
}

func resourceMonitorDiagnosticSettingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("deleting Monitor Diagnostics Setting %q for Resource %q: %+v", id.Name, id.ResourceUri, err)
		}
	}

	// API appears to be eventually consistent (identified during tainting this resource)
	log.Printf("[DEBUG] Waiting for Monitor Diagnostic Setting %q for Resource %q to disappear", id.Name, id.ResourceUri)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   monitorDiagnosticSettingDeletedRefreshFunc(ctx, client, *id),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Monitor Diagnostic Setting %q for Resource %q to become available: %s", id.Name, id.ResourceUri, err)
	}

	return nil
}

func monitorDiagnosticSettingDeletedRefreshFunc(ctx context.Context, client *diagnosticsettings.DiagnosticSettingsClient, targetResourceId diagnosticsettings.ScopedDiagnosticSettingId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, targetResourceId)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}
			return nil, "", fmt.Errorf("issuing read request in monitorDiagnosticSettingDeletedRefreshFunc: %s", err)
		}

		return res, "Exists", nil
	}
}

func expandMonitorDiagnosticsSettingsLogs(input []interface{}) []diagnosticsettings.LogSettings {
	results := make([]diagnosticsettings.LogSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		categoryGroup := v["category_group"].(string)
		enabled := v["enabled"].(bool)
		policiesRaw := v["retention_policy"].([]interface{})
		var retentionPolicy *diagnosticsettings.RetentionPolicy
		if len(policiesRaw) != 0 {
			policyRaw := policiesRaw[0].(map[string]interface{})
			retentionDays := policyRaw["days"].(int)
			retentionEnabled := policyRaw["enabled"].(bool)
			retentionPolicy = &diagnosticsettings.RetentionPolicy{
				Days:    int64(retentionDays),
				Enabled: retentionEnabled,
			}
		}

		output := diagnosticsettings.LogSettings{
			Enabled:         enabled,
			RetentionPolicy: retentionPolicy,
		}
		if category != "" {
			output.Category = utils.String(category)
		} else {
			output.CategoryGroup = utils.String(categoryGroup)
		}

		results = append(results, output)
	}

	return results
}

func flattenMonitorDiagnosticLogs(input *[]diagnosticsettings.LogSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Category != nil {
			output["category"] = *v.Category
		}

		if v.CategoryGroup != nil {
			output["category_group"] = *v.CategoryGroup
		}

		output["enabled"] = v.Enabled

		policies := make([]interface{}, 0)

		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			outputPolicy := make(map[string]interface{})

			outputPolicy["days"] = int(inputPolicy.Days)

			outputPolicy["enabled"] = inputPolicy.Enabled

			policies = append(policies, outputPolicy)
		}

		output["retention_policy"] = policies

		results = append(results, output)
	}

	return results
}

func expandMonitorDiagnosticsSettingsMetrics(input []interface{}) []diagnosticsettings.MetricSettings {
	results := make([]diagnosticsettings.MetricSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		enabled := v["enabled"].(bool)

		policiesRaw := v["retention_policy"].([]interface{})
		var retentionPolicy *diagnosticsettings.RetentionPolicy
		if len(policiesRaw) > 0 && policiesRaw[0] != nil {
			policyRaw := policiesRaw[0].(map[string]interface{})
			retentionDays := policyRaw["days"].(int)
			retentionEnabled := policyRaw["enabled"].(bool)
			retentionPolicy = &diagnosticsettings.RetentionPolicy{
				Days:    int64(retentionDays),
				Enabled: retentionEnabled,
			}
		}
		output := diagnosticsettings.MetricSettings{
			Category:        utils.String(category),
			Enabled:         enabled,
			RetentionPolicy: retentionPolicy,
		}

		results = append(results, output)
	}

	return results
}

func flattenMonitorDiagnosticMetrics(input *[]diagnosticsettings.MetricSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Category != nil {
			output["category"] = *v.Category
		}

		output["enabled"] = v.Enabled

		policies := make([]interface{}, 0)

		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			outputPolicy := make(map[string]interface{})

			outputPolicy["days"] = int(inputPolicy.Days)

			outputPolicy["enabled"] = inputPolicy.Enabled

			policies = append(policies, outputPolicy)
		}

		output["retention_policy"] = policies

		results = append(results, output)
	}

	return results
}

func ParseMonitorDiagnosticId(monitorId string) (*diagnosticsettings.ScopedDiagnosticSettingId, error) {
	v := strings.Split(monitorId, "|")
	if len(v) != 2 {
		return nil, fmt.Errorf("Expected the Monitor Diagnostics ID to be in the format `{resourceId}|{name}` but got %d segments", len(v))
	}

	identifier := diagnosticsettings.ScopedDiagnosticSettingId{
		ResourceUri: v[0],
		Name:        v[1],
	}
	return &identifier, nil
}

func resourceMonitorDiagnosticLogSettingHash(input interface{}) int {
	var buf bytes.Buffer
	if rawData, ok := input.(map[string]interface{}); ok {
		if category, ok := rawData["category"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", category.(string)))
		}
		if categoryGroup, ok := rawData["category_group"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", categoryGroup.(string)))
		}
		if enabled, ok := rawData["enabled"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", enabled.(bool)))
		}
		if policy, ok := rawData["retention_policy"].(map[string]interface{}); ok {
			if policyEnabled, ok := policy["enabled"]; ok {
				buf.WriteString(fmt.Sprintf("%t-", policyEnabled.(bool)))
			}
			if days, ok := policy["days"]; ok {
				buf.WriteString(fmt.Sprintf("%d-", days.(int)))
			}
		}
	}
	return pluginsdk.HashString(buf.String())
}

func resourceMonitorDiagnosticMetricsSettingHash(input interface{}) int {
	var buf bytes.Buffer
	if rawData, ok := input.(map[string]interface{}); ok {
		if category, ok := rawData["category"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", category.(string)))
		}
		if enabled, ok := rawData["enabled"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", enabled.(bool)))
		}
		if policy, ok := rawData["retention_policy"].(map[string]interface{}); ok {
			if policyEnabled, ok := policy["enabled"]; ok {
				buf.WriteString(fmt.Sprintf("%t-", policyEnabled.(bool)))
			}
			if days, ok := policy["days"]; ok {
				buf.WriteString(fmt.Sprintf("%d-", days.(int)))
			}
		}
	}
	return pluginsdk.HashString(buf.String())
}
