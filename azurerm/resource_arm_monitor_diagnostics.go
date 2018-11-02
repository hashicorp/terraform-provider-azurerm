package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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
				// TODO: validation
			},

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"event_hub_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"event_hub_authorization_rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
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

	metrics := expandMetricsConfiguration(*allMetricSettings, disabledSettings, retentionDays)
	logs := expandLogConfiguration(*allLogSettings, disabledSettings, retentionDays)

	// TODO: fix the schema
	properties := insights.DiagnosticSettingsResource{
		DiagnosticSettings: &insights.DiagnosticSettings{
			StorageAccountID:            utils.String(storageAccountId),
			WorkspaceID:                 utils.String(workspaceId),
			EventHubAuthorizationRuleID: utils.String(eventHubAuthorizationRuleId),
			EventHubName:                utils.String(eventHubName),
			Metrics:                     metrics,
			Logs:                        logs,
			// TODO: add to the schema
			//ServiceBusRuleID: utils.String(serviceBusRuleId),
		},
	}

	_, err = client.CreateOrUpdate(ctx, targetResourceId, properties, name)
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

		metricSetting := insights.MetricSettings{
			Category: &settingAsString,
			Enabled:  &enabled,
			RetentionPolicy: &insights.RetentionPolicy{
				Days:    &retentionDays,
				Enabled: &enabled,
			},
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

		logSetting := insights.LogSettings{
			Category: &settingAsString,
			Enabled:  &enabled,
			RetentionPolicy: &insights.RetentionPolicy{
				Days:    &retentionDays,
				Enabled: &enabled,
			},
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

func getAllDiagnosticSettings(targetResourceId string, meta interface{}) (*[]interface{}, *[]interface{}, error) {
	client := meta.(*ArmClient).monitorDiagnosticSettingsCategoryClient
	ctx := meta.(*ArmClient).StopContext
	returnMetricSettings := make([]interface{}, 0)
	returnLogSettings := make([]interface{}, 0)

	categoryList, err := client.List(ctx, targetResourceId)
	if err != nil {
		return nil, nil, err
	}

	for _, item := range *categoryList.Value {
		if item.DiagnosticSettingsCategory.CategoryType == "Metrics" {
			returnMetricSettings = append(returnMetricSettings, *item.Name)
		}
		if item.DiagnosticSettingsCategory.CategoryType == "Logs" {
			returnLogSettings = append(returnLogSettings, *item.Name)
		}
	}

	return &returnMetricSettings, &returnLogSettings, nil
}
