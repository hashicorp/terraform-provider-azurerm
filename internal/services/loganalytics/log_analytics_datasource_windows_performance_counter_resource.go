// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsDataSourceWindowsPerformanceCounter() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogAnalyticsDataSourceWindowsPerformanceCounterCreateUpdate,
		Read:   resourceLogAnalyticsDataSourceWindowsPerformanceCounterRead,
		Update: resourceLogAnalyticsDataSourceWindowsPerformanceCounterCreateUpdate,
		Delete: resourceLogAnalyticsDataSourceWindowsPerformanceCounterDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := datasources.ParseDataSourceID(id)
			return err
		}, importLogAnalyticsDataSource(datasources.DataSourceKindWindowsPerformanceCounter)),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.WindowsPerformanceCounterV0ToV1{},
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

			"counter_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"instance_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"interval_seconds": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(10, math.MaxInt32),
			},

			"object_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

// this should not have been merged, needs to be fixed once https://github.com/Azure/azure-rest-api-specs/issues/9072 has been addressed
type dataSourceWindowsPerformanceCounterProperty struct {
	CounterName     string `json:"counterName"`
	InstanceName    string `json:"instanceName"`
	IntervalSeconds int    `json:"intervalSeconds"`
	ObjectName      string `json:"objectName"`
}

func resourceLogAnalyticsDataSourceWindowsPerformanceCounterCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := datasources.NewDataSourceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("workspace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("failed to check for existing Windows Performance Counter %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_log_analytics_datasource_windows_performance_counter", id.ID())
		}
	}

	params := datasources.DataSource{
		Kind: datasources.DataSourceKindWindowsPerformanceCounter,
		Properties: &dataSourceWindowsPerformanceCounterProperty{
			CounterName:     d.Get("counter_name").(string),
			InstanceName:    d.Get("instance_name").(string),
			IntervalSeconds: d.Get("interval_seconds").(int),
			ObjectName:      d.Get("object_name").(string),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogAnalyticsDataSourceWindowsPerformanceCounterRead(d, meta)
}

func resourceLogAnalyticsDataSourceWindowsPerformanceCounterRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Windows Performance Counter %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Windows Performance Counter %s: %+v", *id, err)
	}

	d.Set("name", id.DataSourceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("workspace_name", id.WorkspaceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			propStr, err := pluginsdk.FlattenJsonToString(props.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("failed to flatten properties map to json: %+v", err)
			}

			prop := &dataSourceWindowsPerformanceCounterProperty{}
			if err := json.Unmarshal([]byte(propStr), &prop); err != nil {
				return fmt.Errorf("failed to decode properties json: %+v", err)
			}

			d.Set("counter_name", prop.CounterName)
			d.Set("instance_name", prop.InstanceName)
			d.Set("interval_seconds", prop.IntervalSeconds)
			d.Set("object_name", prop.ObjectName)
		}
	}

	return nil
}

func resourceLogAnalyticsDataSourceWindowsPerformanceCounterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := datasources.ParseDataSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("failed to delete Windows Performance Counter %s: %+v", *id, err)
	}

	return nil
}
