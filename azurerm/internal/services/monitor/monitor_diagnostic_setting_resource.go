package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	eventhubParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	eventhubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	logAnalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	logAnalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	storageParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorDiagnosticSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorDiagnosticSettingCreateUpdate,
		Read:   resourceMonitorDiagnosticSettingRead,
		Update: resourceMonitorDiagnosticSettingCreateUpdate,
		Delete: resourceMonitorDiagnosticSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MonitorDiagnosticSettingName,
			},

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"eventhub_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubName(),
			},

			"eventhub_authorization_rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.NamespaceAuthorizationRuleID,
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: logAnalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
			},

			"log_analytics_destination_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"Dedicated",
					"AzureDiagnostics", // Not documented in azure API, but some resource has skew. See: https://github.com/Azure/azure-rest-api-specs/issues/9281
				}, false),
			},

			"log": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"retention_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},

									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
						},
					},
				},
			},

			"metric": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"retention_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},

									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceMonitorDiagnosticSettingCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Diagnostic Settings.")

	name := d.Get("name").(string)
	actualResourceId := d.Get("target_resource_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, actualResourceId, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Diagnostic Setting %q for Resource %q: %s", name, actualResourceId, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_diagnostic_setting", *existing.ID)
		}
	}

	logsRaw := d.Get("log").(*schema.Set).List()
	logs := expandMonitorDiagnosticsSettingsLogs(logsRaw)
	metricsRaw := d.Get("metric").(*schema.Set).List()
	metrics := expandMonitorDiagnosticsSettingsMetrics(metricsRaw)

	// if no blocks are specified  the API "creates" but 404's on Read
	if len(logs) == 0 && len(metrics) == 0 {
		return fmt.Errorf("At least one `log` or `metric` block must be specified")
	}

	// also if there's none enabled
	valid := false
	for _, v := range logs {
		if v.Enabled != nil && *v.Enabled {
			valid = true
			break
		}
	}
	if !valid {
		for _, v := range metrics {
			if v.Enabled != nil && *v.Enabled {
				valid = true
				break
			}
		}
	}

	if !valid {
		return fmt.Errorf("At least one `log` or `metric` must be enabled")
	}

	properties := insights.DiagnosticSettingsResource{
		DiagnosticSettings: &insights.DiagnosticSettings{
			Logs:    &logs,
			Metrics: &metrics,
		},
	}

	valid = false
	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		properties.DiagnosticSettings.EventHubAuthorizationRuleID = utils.String(eventHubAuthorizationRuleId)
		properties.DiagnosticSettings.EventHubName = utils.String(eventHubName)
		valid = true
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		properties.DiagnosticSettings.WorkspaceID = utils.String(workspaceId)
		valid = true
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		properties.DiagnosticSettings.StorageAccountID = utils.String(storageAccountId)
		valid = true
	}

	if v := d.Get("log_analytics_destination_type").(string); v != "" {
		if workspaceId != "" {
			properties.DiagnosticSettings.LogAnalyticsDestinationType = &v
		} else {
			return fmt.Errorf("`log_analytics_workspace_id` must be set for `log_analytics_destination_type` to be used")
		}
	}

	if !valid {
		return fmt.Errorf("Either a `eventhub_authorization_rule_id`, `log_analytics_workspace_id` or `storage_account_id` must be set")
	}

	// the Azure SDK prefixes the URI with a `/` such this makes a bad request if we don't trim the `/`
	targetResourceId := strings.TrimPrefix(actualResourceId, "/")
	if _, err := client.CreateOrUpdate(ctx, targetResourceId, properties, name); err != nil {
		return fmt.Errorf("Error creating Monitor Diagnostics Setting %q for Resource %q: %+v", name, actualResourceId, err)
	}

	read, err := client.Get(ctx, targetResourceId, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Monitor Diagnostics %q for Resource ID %q", name, actualResourceId)
	}

	d.SetId(fmt.Sprintf("%s|%s", actualResourceId, name))

	return resourceMonitorDiagnosticSettingRead(d, meta)
}

func resourceMonitorDiagnosticSettingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	actualResourceId := id.ResourceID
	targetResourceId := strings.TrimPrefix(actualResourceId, "/")
	resp, err := client.Get(ctx, targetResourceId, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Monitor Diagnostics Setting %q was not found for Resource %q - removing from state!", id.Name, actualResourceId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Monitor Diagnostics Setting %q for Resource %q: %+v", id.Name, actualResourceId, err)
	}

	d.Set("name", id.Name)
	d.Set("target_resource_id", id.ResourceID)

	d.Set("eventhub_name", resp.EventHubName)
	eventhubAuthorizationRuleId := ""
	if resp.EventHubAuthorizationRuleID != nil && *resp.EventHubAuthorizationRuleID != "" {
		parsedId, err := eventhubParse.NamespaceAuthorizationRuleID(*resp.EventHubAuthorizationRuleID)
		if err != nil {
			return err
		}

		eventhubAuthorizationRuleId = parsedId.ID()
	}
	d.Set("eventhub_authorization_rule_id", eventhubAuthorizationRuleId)

	workspaceId := ""
	if resp.WorkspaceID != nil && *resp.WorkspaceID != "" {
		parsedId, err := logAnalyticsParse.LogAnalyticsWorkspaceID(*resp.WorkspaceID)
		if err != nil {
			return err
		}

		workspaceId = parsedId.ID()
	}
	d.Set("log_analytics_workspace_id", workspaceId)

	storageAccountId := ""
	if resp.StorageAccountID != nil && *resp.StorageAccountID != "" {
		parsedId, err := storageParse.StorageAccountID(*resp.StorageAccountID)
		if err != nil {
			return err
		}

		storageAccountId = parsedId.ID()
	}
	d.Set("storage_account_id", storageAccountId)

	d.Set("log_analytics_destination_type", resp.LogAnalyticsDestinationType)

	if err := d.Set("log", flattenMonitorDiagnosticLogs(resp.Logs)); err != nil {
		return fmt.Errorf("Error setting `log`: %+v", err)
	}

	if err := d.Set("metric", flattenMonitorDiagnosticMetrics(resp.Metrics)); err != nil {
		return fmt.Errorf("Error setting `metric`: %+v", err)
	}

	return nil
}

func resourceMonitorDiagnosticSettingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	targetResourceId := strings.TrimPrefix(id.ResourceID, "/")
	resp, err := client.Delete(ctx, targetResourceId, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Monitor Diagnostics Setting %q for Resource %q: %+v", id.Name, targetResourceId, err)
		}
	}

	// API appears to be eventually consistent (identified during tainting this resource)
	log.Printf("[DEBUG] Waiting for Monitor Diagnostic Setting %q for Resource %q to disappear", id.Name, id.ResourceID)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   monitorDiagnosticSettingDeletedRefreshFunc(ctx, client, targetResourceId, id.Name),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   d.Timeout(schema.TimeoutDelete),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Monitor Diagnostic Setting %q for Resource %q to become available: %s", id.Name, id.ResourceID, err)
	}

	return nil
}

func monitorDiagnosticSettingDeletedRefreshFunc(ctx context.Context, client *insights.DiagnosticSettingsClient, targetResourceId string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, targetResourceId, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}
			return nil, "", fmt.Errorf("Error issuing read request in monitorDiagnosticSettingDeletedRefreshFunc: %s", err)
		}

		return res, "Exists", nil
	}
}

func expandMonitorDiagnosticsSettingsLogs(input []interface{}) []insights.LogSettings {
	results := make([]insights.LogSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		enabled := v["enabled"].(bool)
		policiesRaw := v["retention_policy"].([]interface{})
		var retentionPolicy *insights.RetentionPolicy
		if len(policiesRaw) != 0 {
			policyRaw := policiesRaw[0].(map[string]interface{})
			retentionDays := policyRaw["days"].(int)
			retentionEnabled := policyRaw["enabled"].(bool)
			retentionPolicy = &insights.RetentionPolicy{
				Days:    utils.Int32(int32(retentionDays)),
				Enabled: utils.Bool(retentionEnabled),
			}
		}

		output := insights.LogSettings{
			Category:        utils.String(category),
			Enabled:         utils.Bool(enabled),
			RetentionPolicy: retentionPolicy,
		}

		results = append(results, output)
	}

	return results
}

func flattenMonitorDiagnosticLogs(input *[]insights.LogSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Category != nil {
			output["category"] = *v.Category
		}

		if v.Enabled != nil {
			output["enabled"] = *v.Enabled
		}

		policies := make([]interface{}, 0)

		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			outputPolicy := make(map[string]interface{})

			if inputPolicy.Days != nil {
				outputPolicy["days"] = int(*inputPolicy.Days)
			}

			if inputPolicy.Enabled != nil {
				outputPolicy["enabled"] = *inputPolicy.Enabled
			}

			policies = append(policies, outputPolicy)
		}

		output["retention_policy"] = policies

		results = append(results, output)
	}

	return results
}

func expandMonitorDiagnosticsSettingsMetrics(input []interface{}) []insights.MetricSettings {
	results := make([]insights.MetricSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		enabled := v["enabled"].(bool)

		policiesRaw := v["retention_policy"].([]interface{})
		var retentionPolicy *insights.RetentionPolicy
		if len(policiesRaw) > 0 && policiesRaw[0] != nil {
			policyRaw := policiesRaw[0].(map[string]interface{})
			retentionDays := policyRaw["days"].(int)
			retentionEnabled := policyRaw["enabled"].(bool)
			retentionPolicy = &insights.RetentionPolicy{
				Days:    utils.Int32(int32(retentionDays)),
				Enabled: utils.Bool(retentionEnabled),
			}
		}
		output := insights.MetricSettings{
			Category:        utils.String(category),
			Enabled:         utils.Bool(enabled),
			RetentionPolicy: retentionPolicy,
		}

		results = append(results, output)
	}

	return results
}

func flattenMonitorDiagnosticMetrics(input *[]insights.MetricSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := make(map[string]interface{})

		if v.Category != nil {
			output["category"] = *v.Category
		}

		if v.Enabled != nil {
			output["enabled"] = *v.Enabled
		}

		policies := make([]interface{}, 0)

		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			outputPolicy := make(map[string]interface{})

			if inputPolicy.Days != nil {
				outputPolicy["days"] = int(*inputPolicy.Days)
			}

			if inputPolicy.Enabled != nil {
				outputPolicy["enabled"] = *inputPolicy.Enabled
			}

			policies = append(policies, outputPolicy)
		}

		output["retention_policy"] = policies

		results = append(results, output)
	}

	return results
}

type monitorDiagnosticId struct {
	ResourceID string
	Name       string
}

func ParseMonitorDiagnosticId(monitorId string) (*monitorDiagnosticId, error) {
	v := strings.Split(monitorId, "|")
	if len(v) != 2 {
		return nil, fmt.Errorf("Expected the Monitor Diagnostics ID to be in the format `{resourceId}|{name}` but got %d segments", len(v))
	}

	identifier := monitorDiagnosticId{
		ResourceID: v[0],
		Name:       v[1],
	}
	return &identifier, nil
}
