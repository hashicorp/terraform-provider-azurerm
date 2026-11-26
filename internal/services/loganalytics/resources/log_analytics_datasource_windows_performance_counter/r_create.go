// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_performance_counter

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogAnalyticsDataSourceWindowsPerformanceCounterCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := datasources.NewDataSourceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("workspace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceLogAnalyticsDataSourceWindowsPerformanceCounterRead(d, meta)
}
