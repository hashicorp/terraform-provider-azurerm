// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsDataSourceWindowsEvent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsDataSourceWindowsEventCreateUpdate,
		Read:   resourceLogAnalyticsDataSourceWindowsEventRead,
		Update: resourceLogAnalyticsDataSourceWindowsEventCreateUpdate,
		Delete: resourceLogAnalyticsDataSourceWindowsEventDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := datasources.ParseDataSourceID(id)
			return err
		}, importLogAnalyticsDataSource(datasources.DataSourceKindWindowsEvent)),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsEventV0ToV1{},
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_name": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.LogAnalyticsWorkspaceName,
			},

			"event_log_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_types": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Set:      set.HashStringIgnoreCase,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					// API backend accepts event_types case-insensitively
					ValidateFunc: validation.StringInSlice([]string{"Error", "Warning", "Information"}, false),
				},
			},
		},
	}
}

// this should not have been merged, needs to be fixed once https://github.com/Azure/azure-rest-api-specs/issues/9072 has been addressed
type dataSourceWindowsEvent struct {
	EventLogName string                            `json:"eventLogName"`
	EventTypes   []dataSourceWindowsEventEventType `json:"eventTypes"`
}

// this should not have been merged, needs to be fixed once https://github.com/Azure/azure-rest-api-specs/issues/9072 has been addressed
type dataSourceWindowsEventEventType struct {
	EventType string `json:"eventType"`
}

func resourceLogAnalyticsDataSourceWindowsEventCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := datasources.NewDataSourceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("workspace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_datasource_windows_event", id.ID())
		}
	}

	params := datasources.DataSource{
		Kind: datasources.DataSourceKindWindowsEvent,
		Properties: &dataSourceWindowsEvent{
			EventLogName: d.Get("event_log_name").(string),
			EventTypes:   expandLogAnalyticsDataSourceWindowsEventEventType(d.Get("event_types").(*pluginsdk.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating Windows Event %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsDataSourceWindowsEventRead(d, meta)
}

func resourceLogAnalyticsDataSourceWindowsEventRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datasources.ParseDataSourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Windows Event %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DataSourceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_name", id.WorkspaceName)

	if model := resp.Model; model != nil {
		if props := resp.Model.Properties; props != nil {
			propStr, err := pluginsdk.FlattenJsonToString(props.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("failed to flatten properties map to json: %+v", err)
			}

			prop := dataSourceWindowsEvent{}
			if err := json.Unmarshal([]byte(propStr), &prop); err != nil {
				return fmt.Errorf("failed to decode properties json: %+v", err)
			}

			d.Set("event_log_name", prop.EventLogName)
			d.Set("event_types", flattenLogAnalyticsDataSourceWindowsEventEventType(prop.EventTypes))
		}
	}

	return nil
}

func resourceLogAnalyticsDataSourceWindowsEventDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datasources.ParseDataSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Windows Event %s: %+v", *id, err)
	}

	return nil
}

func flattenLogAnalyticsDataSourceWindowsEventEventType(eventTypes []dataSourceWindowsEventEventType) []interface{} {
	output := make([]interface{}, 0)
	for _, e := range eventTypes {
		// The casing isn't preserved by the API for event types, so we need to normalise it here until
		// https://github.com/Azure/azure-rest-api-specs/issues/18163 is fixed
		output = append(output, e.EventType)
	}
	return output
}

func expandLogAnalyticsDataSourceWindowsEventEventType(input []interface{}) []dataSourceWindowsEventEventType {
	output := []dataSourceWindowsEventEventType{}
	for _, eventType := range input {
		output = append(output, dataSourceWindowsEventEventType{eventType.(string)})
	}
	return output
}
