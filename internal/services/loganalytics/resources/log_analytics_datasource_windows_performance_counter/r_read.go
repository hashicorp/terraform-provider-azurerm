// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_performance_counter

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsDataSourceWindowsPerformanceCounterRead(d *pluginsdk.ResourceData, meta any) error {
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
			propStr, err := pluginsdk.FlattenJsonToString(props.(map[string]any))
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

// this should not have been merged, needs to be fixed once https://github.com/Azure/azure-rest-api-specs/issues/9072 has been addressed
type dataSourceWindowsPerformanceCounterProperty struct {
	CounterName     string `json:"counterName"`
	InstanceName    string `json:"instanceName"`
	IntervalSeconds int    `json:"intervalSeconds"`
	ObjectName      string `json:"objectName"`
}
