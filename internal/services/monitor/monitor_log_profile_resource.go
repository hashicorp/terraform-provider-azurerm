// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2016-03-01/logprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
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

		DeprecationMessage: "Azure Log Profiles will be retired on 30th September 2026 and will be removed in v4.0 of the AzureRM Provider. More information on the deprecation can be found at https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/activity-log?tabs=powershell#legacy-collection-methods",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := logprofiles.ParseLogProfileID(id)
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
				ValidateFunc: commonids.ValidateStorageAccountID,
			},
			"servicebus_rule_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
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

	id := logprofiles.NewLogProfileID(subscriptionId, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_log_profile", id.ID())
		}
	}

	storageAccountID := d.Get("storage_account_id").(string)
	serviceBusRuleID := d.Get("servicebus_rule_id").(string)
	categories := expandLogProfileCategories(d)
	locations := expandLogProfileLocations(d)
	retentionPolicy := expandAzureRmLogProfileRetentionPolicy(d)

	logProfileProperties := logprofiles.LogProfileProperties{
		Categories:      categories,
		Locations:       locations,
		RetentionPolicy: retentionPolicy,
	}

	if storageAccountID != "" {
		logProfileProperties.StorageAccountId = utils.String(storageAccountID)
	}

	if serviceBusRuleID != "" {
		logProfileProperties.ServiceBusRuleId = utils.String(serviceBusRuleID)
	}

	parameters := logprofiles.LogProfileResource{
		Name:       utils.String(id.LogProfileName),
		Properties: logProfileProperties,
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to be provisioned", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Available"},
		Refresh:                   logProfilesCreateRefreshFunc(ctx, client, id),
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

	id, err := logprofiles.ParseLogProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.LogProfileName)

	if model := resp.Model; model != nil {
		props := model.Properties
		d.Set("storage_account_id", props.StorageAccountId)
		d.Set("servicebus_rule_id", props.ServiceBusRuleId)
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

	id, err := logprofiles.ParseLogProfileID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
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

func expandAzureRmLogProfileRetentionPolicy(d *pluginsdk.ResourceData) logprofiles.RetentionPolicy {
	vs := d.Get("retention_policy").([]interface{})
	v := vs[0].(map[string]interface{})

	enabled := v["enabled"].(bool)
	days := v["days"].(int)
	logProfileRetentionPolicy := logprofiles.RetentionPolicy{
		Enabled: enabled,
		Days:    int64(days),
	}

	return logProfileRetentionPolicy
}

func flattenAzureRmLogProfileLocations(input []string) []string {
	result := make([]string, 0)
	for _, location := range input {
		result = append(result, azure.NormalizeLocation(location))
	}
	return result
}

func flattenAzureRmLogProfileRetentionPolicy(input logprofiles.RetentionPolicy) []interface{} {
	result := make(map[string]interface{})
	result["enabled"] = input.Enabled

	result["days"] = input.Days

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

func logProfilesCreateRefreshFunc(ctx context.Context, client *logprofiles.LogProfilesClient, id logprofiles.LogProfileId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logProfile, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(logProfile.HttpResponse) {
				return nil, "NotFound", nil
			}
			return nil, "", fmt.Errorf("issuing read request in logProfilesCreateRefreshFunc for  %s: %s", id, err)
		}
		return "Available", "Available", nil
	}
}
