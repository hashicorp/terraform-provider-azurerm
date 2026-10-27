// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_event

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

func resourceLogAnalyticsDataSourceWindowsEventRead(d *pluginsdk.ResourceData, meta any) error {
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
			propStr, err := pluginsdk.FlattenJsonToString(props.(map[string]any))
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

func flattenLogAnalyticsDataSourceWindowsEventEventType(eventTypes []dataSourceWindowsEventEventType) []any {
	output := make([]any, 0)
	for _, e := range eventTypes {
		// The casing isn't preserved by the API for event types, so we need to normalise it here until
		// https://github.com/Azure/azure-rest-api-specs/issues/18163 is fixed
		output = append(output, e.EventType)
	}
	return output
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
