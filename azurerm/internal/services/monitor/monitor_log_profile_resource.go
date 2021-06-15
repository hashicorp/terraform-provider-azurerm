package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorLogProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogProfileCreateUpdate,
		Read:   resourceLogProfileRead,
		Update: resourceLogProfileCreateUpdate,
		Delete: resourceLogProfileDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"servicebus_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"locations": {
				Type:     pluginsdk.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:             pluginsdk.TypeString,
					StateFunc:        location.StateFunc,
					DiffSuppressFunc: location.DiffSuppressFunc,
				},
				Set: pluginsdk.HashString,
			},
			"categories": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:             pluginsdk.TypeString,
					DiffSuppressFunc: suppress.CaseDifference,
				},
				Set: pluginsdk.HashString,
			},
			"retention_policy": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"days": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceLogProfileCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.LogProfilesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Log Profile %q: %s", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_log_profile", *existing.ID)
		}
	}

	storageAccountID := d.Get("storage_account_id").(string)
	serviceBusRuleID := d.Get("servicebus_rule_id").(string)
	categories := expandLogProfileCategories(d)
	locations := expandLogProfileLocations(d)
	retentionPolicy := expandAzureRmLogProfileRetentionPolicy(d)

	logProfileProperties := &insights.LogProfileProperties{
		Categories:      &categories,
		Locations:       &locations,
		RetentionPolicy: &retentionPolicy,
	}

	if storageAccountID != "" {
		logProfileProperties.StorageAccountID = utils.String(storageAccountID)
	}

	if serviceBusRuleID != "" {
		logProfileProperties.ServiceBusRuleID = utils.String(serviceBusRuleID)
	}

	parameters := insights.LogProfileResource{
		Name:                 utils.String(name),
		LogProfileProperties: logProfileProperties,
	}

	if _, err := client.CreateOrUpdate(ctx, name, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Log Profile %q: %+v", name, err)
	}

	log.Printf("[DEBUG] Waiting for Log Profile %q to be provisioned", name)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Available"},
		Refresh:                   logProfilesCreateRefreshFunc(ctx, client, name),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for Log Profile %q to become available: %s", name, err)
	}

	read, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Log Profile %q: %+v", name, err)
	}

	d.SetId(*read.ID)

	return resourceLogProfileRead(d, meta)
}

func resourceLogProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.LogProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := ParseLogProfileNameFromID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing log profile name from ID %s: %s", d.Id(), err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Log Profile %q was not found - removing from state!", name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Log Profile %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	if props := resp.LogProfileProperties; props != nil {
		d.Set("storage_account_id", props.StorageAccountID)
		d.Set("servicebus_rule_id", props.ServiceBusRuleID)
		d.Set("categories", props.Categories)

		if err := d.Set("locations", flattenAzureRmLogProfileLocations(props.Locations)); err != nil {
			return fmt.Errorf("Error setting `locations`: %+v", err)
		}

		if err := d.Set("retention_policy", flattenAzureRmLogProfileRetentionPolicy(props.RetentionPolicy)); err != nil {
			return fmt.Errorf("Error setting `retention_policy`: %+v", err)
		}
	}

	return nil
}

func resourceLogProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.LogProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := ParseLogProfileNameFromID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing log profile name from ID %s: %s", d.Id(), err)
	}

	if _, err = client.Delete(ctx, name); err != nil {
		return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
	}

	return nil
}

func expandLogProfileCategories(d *pluginsdk.ResourceData) []string {
	logProfileCategories := d.Get("categories").(*pluginsdk.Set).List()
	categories := make([]string, 0)

	for _, category := range logProfileCategories {
		categories = append(categories, category.(string))
	}

	return categories
}

func expandLogProfileLocations(d *pluginsdk.ResourceData) []string {
	logProfileLocations := d.Get("locations").(*pluginsdk.Set).List()
	locations := make([]string, 0)

	for _, location := range logProfileLocations {
		locations = append(locations, azure.NormalizeLocation(location.(string)))
	}

	return locations
}

func expandAzureRmLogProfileRetentionPolicy(d *pluginsdk.ResourceData) insights.RetentionPolicy {
	vs := d.Get("retention_policy").([]interface{})
	v := vs[0].(map[string]interface{})

	enabled := v["enabled"].(bool)
	days := v["days"].(int)
	logProfileRetentionPolicy := insights.RetentionPolicy{
		Enabled: utils.Bool(enabled),
		Days:    utils.Int32(int32(days)),
	}

	return logProfileRetentionPolicy
}

func flattenAzureRmLogProfileLocations(input *[]string) []string {
	result := make([]string, 0)
	if input != nil {
		for _, location := range *input {
			result = append(result, azure.NormalizeLocation(location))
		}
	}

	return result
}

func flattenAzureRmLogProfileRetentionPolicy(input *insights.RetentionPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	if input.Enabled != nil {
		result["enabled"] = *input.Enabled
	}

	if input.Days != nil {
		result["days"] = *input.Days
	}

	return []interface{}{result}
}

func ParseLogProfileNameFromID(id string) (string, error) {
	components := strings.Split(id, "/")

	if len(components) == 0 {
		return "", fmt.Errorf("Azure Log Profile ID is empty or not formatted correctly: %s", id)
	}

	if len(components) != 7 {
		return "", fmt.Errorf("Azure Log Profile ID should have 6 segments, got %d: '%s'", len(components)-1, id)
	}

	return components[6], nil
}

func logProfilesCreateRefreshFunc(ctx context.Context, client *insights.LogProfilesClient, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logProfile, err := client.Get(ctx, name)
		if err != nil {
			if utils.ResponseWasNotFound(logProfile.Response) {
				return nil, "NotFound", nil
			}
			return nil, "", fmt.Errorf("Error issuing read request in logProfilesCreateRefreshFunc for Log profile %q: %s", name, err)
		}
		return "Available", "Available", nil
	}
}
