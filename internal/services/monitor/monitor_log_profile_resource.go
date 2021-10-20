package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorLogProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogProfileCreateUpdate,
		Read:   resourceLogProfileRead,
		Update: resourceLogProfileCreateUpdate,
		Delete: resourceLogProfileDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LogProfileID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LogProfileUpgradeV0ToV1{},
		}),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLogProfileID(subscriptionId, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
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
		Name:                 utils.String(id.Name),
		LogProfileProperties: logProfileProperties,
	}

	if _, err := client.CreateOrUpdate(ctx, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for Log Profile %q to be provisioned", id.Name)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Available"},
		Refresh:                   logProfilesCreateRefreshFunc(ctx, client, id.Name),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 5,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogProfileRead(d, meta)
}

func resourceLogProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.LogProfilesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Log Profile %q was not found - removing from state!", id.Name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	if props := resp.LogProfileProperties; props != nil {
		d.Set("storage_account_id", props.StorageAccountID)
		d.Set("servicebus_rule_id", props.ServiceBusRuleID)
		d.Set("categories", props.Categories)

		if err := d.Set("locations", flattenAzureRmLogProfileLocations(props.Locations)); err != nil {
			return fmt.Errorf("setting `locations`: %+v", err)
		}

		if err := d.Set("retention_policy", flattenAzureRmLogProfileRetentionPolicy(props.RetentionPolicy)); err != nil {
			return fmt.Errorf("setting `retention_policy`: %+v", err)
		}
	}

	return nil
}

func resourceLogProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.LogProfilesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LogProfileID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
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
			return nil, "", fmt.Errorf("issuing read request in logProfilesCreateRefreshFunc for Log profile %q: %s", name, err)
		}
		return "Available", "Available", nil
	}
}
