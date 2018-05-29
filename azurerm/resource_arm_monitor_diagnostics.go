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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"storage_account_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
	metricSettings := d.Get("metric_settings").(*schema.Set)
	logSettings := d.Get("log_settings").(*schema.Set)

	diagnosticSettings := &insights.DiagnosticSettings{}

	if metricSettings != nil {
		diagnosticSettings.Metrics = expandMetricsConfiguration(metricSettings)
	}

	if logSettings != nil {
		diagnosticSettings.Logs = expandLogConfiguration(logSettings)
	}

	if len(storageAccountId) > 0 {
		diagnosticSettings.StorageAccountID = &storageAccountId
	}

	if len(workspaceId) > 0 {
		diagnosticSettings.WorkspaceID = &workspaceId
	}

	if len(eventHubAuthorizationRuleId) > 0 && len(eventHubName) > 0 {
		diagnosticSettings.EventHubAuthorizationRuleID = &eventHubAuthorizationRuleId
		diagnosticSettings.EventHubName = &eventHubName
	}

	_, err := client.CreateOrUpdate(
		ctx,
		resourceId,
		insights.DiagnosticSettingsResource{
			Name:               &name,
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

	monitoringId := parseMonitorDiagnosticId(d.Id())

	resp, err := client.Get(ctx, monitoringId.ResourceID, monitoringId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Diagnostic Setting %s: %+v", monitoringId.Name, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", *resp.Name)

	// ID of base resource is not returned by API, so we have to guess here
	monitoringId = parseMonitorDiagnosticId(d.Id())
	d.Set("resource_id", monitoringId.ResourceID)

	if resp.StorageAccountID != nil {
		d.Set("storage_account_id", *resp.StorageAccountID)
	}

	if resp.EventHubName != nil {
		d.Set("event_hub_name", *resp.EventHubName)
	}

	if resp.EventHubAuthorizationRuleID != nil {
		d.Set("event_hub_authorization_rule_id", *resp.EventHubAuthorizationRuleID)
	}

	if resp.WorkspaceID != nil {
		d.Set("workspace_id", *resp.WorkspaceID)
	}

	d.Set("metric_settings", flattenMetricsConfiguration(*resp.Metrics))
	d.Set("log_settings", flattenLogConfiguration(*resp.Logs))

	return nil
}

func flattenMetricsConfiguration(metricSettings []insights.MetricSettings) []interface{} {
	returnConfiguration := make([]interface{}, 0, len(metricSettings))

	if metricSettings == nil {
		return returnConfiguration
	}

	for _, setting := range metricSettings {
		currentSetting := make(map[string]interface{})

		currentSetting["category"] = *setting.Category
		currentSetting["retention_days"] = *setting.RetentionPolicy.Days
		if setting.TimeGrain != nil {
			currentSetting["time_grain"] = *setting.TimeGrain
		}

		returnConfiguration = append(returnConfiguration, currentSetting)
	}

	return returnConfiguration
}

func flattenLogConfiguration(logSettings []insights.LogSettings) []interface{} {
	returnConfiguration := make([]interface{}, 0, len(logSettings))

	if logSettings == nil {
		return returnConfiguration
	}

	for _, setting := range logSettings {
		currentSetting := make(map[string]interface{})

		currentSetting["category"] = *setting.Category
		currentSetting["retention_days"] = *setting.RetentionPolicy.Days

		returnConfiguration = append(returnConfiguration, currentSetting)
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
