package loganalytics

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLogAnalyticsDataSourceWindowsEvent() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsDataSourceWindowsEventCreateUpdate,
		Read:   resourceLogAnalyticsDataSourceWindowsEventRead,
		Update: resourceLogAnalyticsDataSourceWindowsEventCreateUpdate,
		Delete: resourceLogAnalyticsDataSourceWindowsEventDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DataSourceID(id)
			return err
		}, importLogAnalyticsDataSource(operationalinsights.WindowsEvent)),

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

			"resource_group_name": azure.SchemaResourceGroupName(),

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
					ValidateFunc:     validation.StringInSlice([]string{"error", "warning", "information"}, !features.ThreePointOh()),
					DiffSuppressFunc: suppress.CaseDifferenceV2Only,
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

	id := parse.NewDataSourceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("workspace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_log_analytics_datasource_windows_event", id.ID())
		}
	}

	params := operationalinsights.DataSource{
		Kind: operationalinsights.WindowsEvent,
		Properties: &dataSourceWindowsEvent{
			EventLogName: d.Get("event_log_name").(string),
			EventTypes:   expandLogAnalyticsDataSourceWindowsEventEventType(d.Get("event_types").(*pluginsdk.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, params); err != nil {
		return fmt.Errorf("creating Windows Event %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsDataSourceWindowsEventRead(d, meta)
}

func resourceLogAnalyticsDataSourceWindowsEventRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Windows Event %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("workspace_name", id.WorkspaceName)
	if props := resp.Properties; props != nil {
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

	return nil
}

func resourceLogAnalyticsDataSourceWindowsEventDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Windows Event %s: %+v", *id, err)
	}

	return nil
}

func flattenLogAnalyticsDataSourceWindowsEventEventType(eventTypes []dataSourceWindowsEventEventType) []interface{} {
	output := make([]interface{}, 0)
	for _, e := range eventTypes {
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
