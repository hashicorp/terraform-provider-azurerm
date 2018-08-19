package azurerm

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogProfileCreateOrUpdate,
		Read:   resourceArmLogProfileRead,
		Update: resourceArmLogProfileCreateOrUpdate,
		Delete: resourceArmLogProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"service_bus_rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},
			"locations": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					StateFunc:        azureRMNormalizeLocation,
					DiffSuppressFunc: azureRMSuppressLocationDiff,
				},
			},
			"categories": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				},
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
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceArmLogProfileCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logProfilesClient
	ctx := meta.(*ArmClient).StopContext

	categories := expandLogProfileCategories(d)
	locations := expandLogProfileLocations(d)
	retentionPolicy, err := expandAzureRmLogProfileRetentionPolicy(d)
	if err != nil {
		return err
	}

	logProfileProperties := insights.LogProfileProperties{
		Categories:      &categories,
		Locations:       &locations,
		RetentionPolicy: retentionPolicy,
	}

	storageAccountID := d.Get("storage_account_id").(string)
	if storageAccountID != "" {
		logProfileProperties.StorageAccountID = utils.String(storageAccountID)
	}

	serviceBusRuleID := d.Get("service_bus_rule_id").(string)
	if serviceBusRuleID != "" {
		logProfileProperties.ServiceBusRuleID = utils.String(serviceBusRuleID)
	}

	name := d.Get("name").(string)
	parameters := insights.LogProfileResource{
		Name:                 utils.String(name),
		LogProfileProperties: &logProfileProperties,
	}

	result, createErr := client.CreateOrUpdate(ctx, name, parameters)
	if createErr != nil {
		return fmt.Errorf("Error Creating/Updating Log Profile %q: %+v", name, createErr)
	}

	// Wait for Log Profile to become available
	err = resource.Retry(300*time.Second, retryLogProfilesClientGet(name, meta))
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, name)
	if err != nil {
		return err
	}

	d.SetId(*read.ID)

	return resourceArmLogProfileRead(d, meta)
}

func resourceArmLogProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logProfilesClient
	ctx := meta.(*ArmClient).StopContext

	name, err := parseLogProfileNameFromID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Log Profile %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	if props := resp.LogProfileProperties; props != nil {
		if props.StorageAccountID != nil {
			d.Set("storage_account_id", props.StorageAccountID)
		} else {
			d.Set("storage_account_id", "")
		}

		if props.ServiceBusRuleID != nil {
			d.Set("service_bus_rule_id", props.ServiceBusRuleID)
		} else {
			d.Set("service_bus_rule_id", "")
		}

		d.Set("locations", props.Locations)
		d.Set("categories", props.Categories)

		d.Set("retention_policy", flattenAzureRmLogProfileRetentionPolicy(props.RetentionPolicy))
	}

	return nil
}

func resourceArmLogProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).logProfilesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)

	_, err := client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting Log Profile %q: %+v", name, err)
	}

	return nil
}

func expandLogProfileCategories(d *schema.ResourceData) []string {
	logProfileCategories := d.Get("categories").([]interface{})
	categories := []string{}

	for _, category := range logProfileCategories {
		categories = append(categories, category.(string))
	}

	return categories
}

func expandLogProfileLocations(d *schema.ResourceData) []string {
	logProfileLocations := d.Get("locations").([]interface{})
	locations := []string{}

	for _, location := range logProfileLocations {
		locations = append(locations, location.(string))
	}

	return locations
}

func expandAzureRmLogProfileRetentionPolicy(d *schema.ResourceData) (*insights.RetentionPolicy, error) {
	retentionPolicies := d.Get("retention_policy").([]interface{})

	if len(retentionPolicies) > 0 {
		retentionPolicy := retentionPolicies[0].(map[string]interface{})
		logProfileRetentionPolicy := &insights.RetentionPolicy{
			Enabled: utils.Bool(retentionPolicy["enabled"].(bool)),
			Days:    utils.Int32(int32(retentionPolicy["days"].(int))),
		}

		return logProfileRetentionPolicy, nil
	}

	return nil, fmt.Errorf("[ERROR] Retention policy must be set")
}

func flattenAzureRmLogProfileRetentionPolicy(input *insights.RetentionPolicy) []interface{} {
	result := make(map[string]interface{})
	result["enabled"] = *input.Enabled
	result["days"] = *input.Days

	return []interface{}{result}
}

func retryLogProfilesClientGet(name string, meta interface{}) func() *resource.RetryError {
	return func() *resource.RetryError {
		client := meta.(*ArmClient).logProfilesClient
		ctx := meta.(*ArmClient).StopContext

		read, err := client.Get(ctx, name)
		if err != nil {
			return resource.RetryableError(err)
		}

		return nil
	}
}

func parseLogProfileNameFromID(id string) (string, error) {
	components := strings.Split(id, "/")

	if len(components) == 0 {
		return "", fmt.Errorf("Azure Log Profile ID is empty or not formatted correctly: %s", id)
	}

	if len(components) != 7 {
		return "", fmt.Errorf("Azure Log Profile ID should have 6 segments, got %d: '%s'", len(components)-1, id)
	}

	return components[6], nil
}
