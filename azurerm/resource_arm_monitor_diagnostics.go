package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorDiagnostics() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorDiagnosticsCreate,
		Read:   resourceArmMonitorDiagnosticsRead,
		Update: resourceArmMonitorDiagnosticsCreate,
		Delete: resourceArmMonitorDiagnosticsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"event_hub_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"event_hub_authorization_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metric_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"time_grain": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"log_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 16,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmMonitorDiagnosticsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorDiagnosticSettingsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Diagnostic Settings.")

	name := d.Get("name").(string)
	resourceId := d.Get("resource_id").(string)
	storageAccountId := d.Get("storage_account_id").(string)
	eventHubName := d.Get("event_hub_name").(string)
	eventHubAuthorizationRuleId := d.Get("event_hub_authorization_rule_id").(string)
	workspaceId := d.Get("workspace_id").(string)
	metricSettings := d.Get("metric_settings")
	logSettings := d.Get("log_settings")

	diagnosticSettings := &insights.DiagnosticSettings{}

	if metricSettings != nil {
		diagnosticSettings.Metrics = expandMetricsConfiguration(metricSettings.(*schema.Set))
	}

	if logSettings != nil {
		diagnosticSettings.Logs = expandLogConfiguration(logSettings.(*schema.Set))
	}

	if len(storageAccountId) > 0 {
		diagnosticSettings.StorageAccountID = utils.String(storageAccountId)
	}

	if len(workspaceId) > 0 {
		diagnosticSettings.WorkspaceID = utils.String(workspaceId)
	}

	if len(eventHubAuthorizationRuleId) > 0 && len(eventHubName) > 0 {
		diagnosticSettings.EventHubAuthorizationRuleID = utils.String(eventHubAuthorizationRuleId)
		diagnosticSettings.EventHubName = utils.String(eventHubName)
	}

	_, err := client.CreateOrUpdate(
		ctx,
		resourceId,
		insights.DiagnosticSettingsResource{
			Name:               utils.String(name),
			DiagnosticSettings: diagnosticSettings,
		},
		name)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceId, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Diagnostic Settings")
	}

	d.SetId(*read.ID)

	return resourceArmMonitorDiagnosticsRead(d, meta)
}

func resourceArmMonitorDiagnosticsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorDiagnosticSettingsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resource_id := d.Get("resource_id").(string)

	resp, err := client.Get(ctx, resource_id, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure KeyVault %s: %+v", name, err)
	}
	d.SetId(*resp.ID)

	if resp.DiagnosticSettings.StorageAccountID != nil {
		d.Set("storage_account_id", *resp.DiagnosticSettings.StorageAccountID)
	}

	if resp.DiagnosticSettings.EventHubName != nil {
		d.Set("event_hub_name", *resp.DiagnosticSettings.EventHubName)
	}

	if resp.DiagnosticSettings.EventHubAuthorizationRuleID != nil {
		d.Set("event_hub_authorization_rule_id", *resp.DiagnosticSettings.EventHubAuthorizationRuleID)
	}

	if resp.DiagnosticSettings.WorkspaceID != nil {
		d.Set("workspace_id", *resp.DiagnosticSettings.WorkspaceID)
	}

	d.Set("metric_settings", flattenMetricsConfiguration(*resp.DiagnosticSettings.Metrics))
	d.Set("log_settings", flattenLogConfiguration(*resp.DiagnosticSettings.Logs))

	return nil
}

func flattenMetricsConfiguration(metricsSettings []insights.MetricSettings) []interface{} {
	returnConfiguration := make([]interface{}, 0, len(metricsSettings))

	if metricsSettings == nil {
		return returnConfiguration
	}

	for _, setting := range metricsSettings {
		metricSetting := make(map[string]interface{})

		metricSetting["category"] = setting.Category
		metricSetting["retention_days"] = setting.RetentionPolicy.Days
		metricSetting["time_grain"] = setting.TimeGrain

		returnConfiguration = append(returnConfiguration, metricSetting)
	}

	return returnConfiguration
}

func flattenLogConfiguration(logSettings []insights.LogSettings) []interface{} {
	returnConfiguration := make([]interface{}, 0, len(logSettings))

	if logSettings == nil {
		return returnConfiguration
	}

	for _, setting := range logSettings {
		logSetting := make(map[string]interface{})

		logSetting["category"] = setting.Category
		logSetting["retention_days"] = setting.RetentionPolicy.Days

		returnConfiguration = append(returnConfiguration, logSetting)
	}

	return returnConfiguration
}

func resourceArmMonitorDiagnosticsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorDiagnosticSettingsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resource_id := d.Get("resource_id").(string)

	_, err := client.Delete(ctx, resource_id, name)

	return err
}

func expandMetricsConfiguration(metricsSettings *schema.Set) *[]insights.MetricSettings {
	returnMetricsSettings := make([]insights.MetricSettings, 0, metricsSettings.Len())

	for _, setting := range metricsSettings.List() {
		settingsMap := setting.(map[string]interface{})
		category := settingsMap["category"].(string)
		retentionDays := int32(settingsMap["retention_days"].(int))
		timeGrain := settingsMap["time_grain"].(string)
		enabled := true

		retentionPolicy := insights.RetentionPolicy{
			Days:    &retentionDays,
			Enabled: &enabled,
		}

		metricSetting := insights.MetricSettings{
			Category:        &category,
			Enabled:         &enabled,
			RetentionPolicy: &retentionPolicy,
			TimeGrain:       &timeGrain,
		}
		returnMetricsSettings = append(returnMetricsSettings, metricSetting)
	}
	return &returnMetricsSettings
}

func expandLogConfiguration(logSettings *schema.Set) *[]insights.LogSettings {
	returnLogSettings := make([]insights.LogSettings, 0, logSettings.Len())

	for _, setting := range logSettings.List() {
		settingsMap := setting.(map[string]interface{})
		category := settingsMap["category"].(string)
		retentionDays := int32(settingsMap["retention_days"].(int))
		enabled := true

		retentionPolicy := insights.RetentionPolicy{
			Days:    &retentionDays,
			Enabled: &enabled,
		}

		logSetting := insights.LogSettings{
			Category:        &category,
			Enabled:         &enabled,
			RetentionPolicy: &retentionPolicy,
		}
		returnLogSettings = append(returnLogSettings, logSetting)
	}
	return &returnLogSettings
}
