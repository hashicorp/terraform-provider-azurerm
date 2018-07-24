package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type monitorDiagnosticId struct {
	resourceID string
	name       string
}

func resourceArmMonitorDiagnostics() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorDiagnosticsCreateUpdate,
		Read:   resourceArmMonitorDiagnosticsRead,
		Update: resourceArmMonitorDiagnosticsCreateUpdate,
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

			"target_resource_id": {
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

			"disabled_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceArmMonitorDiagnosticsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorDiagnosticSettingsClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Diagnostic Settings.")

	name := d.Get("name").(string)
	targetResourceId := d.Get("target_resource_id").(string)
	storageAccountId := d.Get("storage_account_id").(string)
	eventHubName := d.Get("event_hub_name").(string)
	eventHubAuthorizationRuleId := d.Get("event_hub_authorization_rule_id").(string)
	workspaceId := d.Get("workspace_id").(string)
	disabledSettings := d.Get("disabled_settings").([]interface{})
	retentionDays := d.Get("retention_days").(int)

	// TODO: I think this wants to be a Data Source?
	allMetricSettings, allLogSettings, err := getAllDiagnosticSettings(targetResourceId, meta)
	if err != nil {
		return err
	}

	allSettings := append(*allMetricSettings, *allLogSettings...)
	if !utils.StringSliceContainsStringSlice(disabledSettings, allSettings) {
		return fmt.Errorf("Invalid value for disabled settings provided, use one or multiple of: %q", allSettings)
	}

	// TODO: remove this
	if len(allSettings) == len(disabledSettings) {
		return fmt.Errorf("You can not disable all settings, rather delete diagnostic logging")
	}

	// TODO: fix the schema
	diagnosticSettings := &insights.DiagnosticSettings{}
	diagnosticSettings.Metrics = expandMetricsConfiguration(*allMetricSettings, disabledSettings, retentionDays)
	diagnosticSettings.Logs = expandLogConfiguration(*allLogSettings, disabledSettings, retentionDays)

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

	resource := insights.DiagnosticSettingsResource{
		DiagnosticSettings: diagnosticSettings,
	}

	_, err = client.CreateOrUpdate(ctx, targetResourceId, resource, name)
	if err != nil {
		return fmt.Errorf("Error creating Diagnostics Setting %q (Resource ID %q): %+v", name, targetResourceId, err)
	}

	read, err := client.Get(ctx, targetResourceId, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Monitor Diagnostics %q for Resource ID %q", name, targetResourceId)
	}

	d.SetId(fmt.Sprintf("%s|%s", targetResourceId, name))

	return resourceArmMonitorDiagnosticsRead(d, meta)
}

func resourceArmMonitorDiagnosticsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorDiagnosticSettingsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.resourceID, id.name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Monitor Diagnostics Setting %q was not found for Resource ID %q - removing from state!", id.name, id.resourceID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Monitor Diagnostics Setting %q for Resource ID %q: %+v", id.name, id.resourceID, err)
	}

	d.Set("name", id.name)
	d.Set("target_resource_id", id.resourceID)
	d.Set("storage_account_id", resp.StorageAccountID)
	d.Set("event_hub_name", resp.EventHubName)
	d.Set("event_hub_authorization_rule_id", resp.EventHubAuthorizationRuleID)
	d.Set("workspace_id", resp.WorkspaceID)

	// TODO: handle crashes/set errors here
	d.Set("disabled_settings", flattenDisabledSettings(*resp.Metrics, *resp.Logs))
	d.Set("retention_days", flattenRetentionDays(*resp.Metrics, *resp.Logs))

	return nil
}

func resourceArmMonitorDiagnosticsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitorDiagnosticSettingsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseMonitorDiagnosticId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.resourceID, id.name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Monitor Diagnostics Setting %q for Resource Id %q: %+v", id.name, id.resourceID, err)
		}
	}

	return nil
}

func expandMetricsConfiguration(allMetricSettings, disabledSettings []interface{}, retentionDays int) *[]insights.MetricSettings {
	returnMetricsSettings := make([]insights.MetricSettings, 0)

	for _, setting := range allMetricSettings {
		settingAsString := setting.(string)
		enabled := true
		if utils.SliceContainsString(disabledSettings, settingAsString) {
			enabled = false
		}
		retentionDays := int32(retentionDays)

		retentionPolicy := insights.RetentionPolicy{
			Days:    &retentionDays,
			Enabled: &enabled,
		}

		metricSetting := insights.MetricSettings{
			Category:        &settingAsString,
			Enabled:         &enabled,
			RetentionPolicy: &retentionPolicy,
		}
		returnMetricsSettings = append(returnMetricsSettings, metricSetting)
	}
	return &returnMetricsSettings
}

func expandLogConfiguration(allLogSettings, disabledSettings []interface{}, retentionDays int) *[]insights.LogSettings {
	returnLogSettings := make([]insights.LogSettings, 0)

	for _, setting := range allLogSettings {
		settingAsString := setting.(string)
		enabled := true
		if utils.SliceContainsString(disabledSettings, settingAsString) {
			enabled = false
		}
		retentionDays := int32(retentionDays)

		retentionPolicy := insights.RetentionPolicy{
			Days:    &retentionDays,
			Enabled: &enabled,
		}

		logSetting := insights.LogSettings{
			Category:        &settingAsString,
			Enabled:         &enabled,
			RetentionPolicy: &retentionPolicy,
		}
		returnLogSettings = append(returnLogSettings, logSetting)
	}
	return &returnLogSettings
}

func flattenDisabledSettings(metricSettings []insights.MetricSettings, logSettings []insights.LogSettings) []interface{} {
	disabledSettings := make([]interface{}, 0)

	for _, setting := range metricSettings {
		category := *setting.Category
		if !*setting.Enabled {
			disabledSettings = append(disabledSettings, category)
		}
	}

	for _, setting := range logSettings {
		category := *setting.Category
		if !*setting.Enabled {
			disabledSettings = append(disabledSettings, category)
		}
	}
	return disabledSettings
}

func flattenRetentionDays(metricSettings []insights.MetricSettings, logSettings []insights.LogSettings) int32 {
	returnSetting := int32(0)

	for _, setting := range metricSettings {
		if *setting.Enabled {
			returnSetting = *setting.RetentionPolicy.Days
		}
	}

	for _, setting := range logSettings {
		if *setting.Enabled {
			returnSetting = *setting.RetentionPolicy.Days
		}
	}

	return returnSetting
}

func parseMonitorDiagnosticId(monitorId string) (*monitorDiagnosticId, error) {
	v := strings.Split(monitorId, "|")
	if len(v) != 2 {
		return nil, fmt.Errorf("Expected the Monitor Diagnostics ID to be in the format `{resourceId}|{name}` but got %d segments", len(v))
	}

	identifier := monitorDiagnosticId{
		resourceID: v[0],
		name:       v[1],
	}
	return &identifier, nil
}
