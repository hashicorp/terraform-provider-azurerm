package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorLogProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogProfileCreateUpdate,
		Read:   resourceLogProfileRead,
		Update: resourceLogProfileCreateUpdate,
		Delete: resourceLogProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"servicebus_rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"locations": {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					StateFunc:        location.StateFunc,
					DiffSuppressFunc: location.DiffSuppressFunc,
				},
				Set: schema.HashString,
			},
			"categories": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: suppress.CaseDifference,
				},
				Set: schema.HashString,
			},
			"retention_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"days": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceLogProfileCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Available"},
		Refresh:                   logProfilesCreateRefreshFunc(ctx, client, name),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Log Profile %q to become available: %s", name, err)
	}

	read, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Log Profile %q: %+v", name, err)
	}

	d.SetId(*read.ID)

	return resourceLogProfileRead(d, meta)
}

func resourceLogProfileRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceLogProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.LogProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := ParseLogProfileNameFromID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing log profile name from ID %s: %s", d.Id(), err)
	}

	_, err = client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
	}

	return nil
}

func expandLogProfileCategories(d *schema.ResourceData) []string {
	logProfileCategories := d.Get("categories").(*schema.Set).List()
	categories := make([]string, 0)

	for _, category := range logProfileCategories {
		categories = append(categories, category.(string))
	}

	return categories
}

func expandLogProfileLocations(d *schema.ResourceData) []string {
	logProfileLocations := d.Get("locations").(*schema.Set).List()
	locations := make([]string, 0)

	for _, location := range logProfileLocations {
		locations = append(locations, azure.NormalizeLocation(location.(string)))
	}

	return locations
}

func expandAzureRmLogProfileRetentionPolicy(d *schema.ResourceData) insights.RetentionPolicy {
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

func logProfilesCreateRefreshFunc(ctx context.Context, client *insights.LogProfilesClient, name string) resource.StateRefreshFunc {
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
