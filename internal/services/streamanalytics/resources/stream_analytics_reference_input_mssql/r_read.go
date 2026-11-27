// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_reference_input_mssql

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsReferenceInputMsSqlRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			reference, ok := props.(inputs.ReferenceInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to Reference Input", *id)
			}

			referenceInputAzureSql, ok := reference.Datasource.(inputs.AzureSqlReferenceInputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to an Azure Sql Reference Input", *id)
			}

			if referenceInputAzureSql.Properties != nil {
				d.Set("server", pointer.From(referenceInputAzureSql.Properties.Server))

				d.Set("database", pointer.From(referenceInputAzureSql.Properties.Database))

				d.Set("username", pointer.From(referenceInputAzureSql.Properties.User))

				refreshType := ""
				if v := referenceInputAzureSql.Properties.RefreshType; v != nil {
					refreshType = string(*v)
				}
				d.Set("refresh_type", refreshType)

				d.Set("refresh_interval_duration", pointer.From(referenceInputAzureSql.Properties.RefreshRate))

				d.Set("full_snapshot_query", pointer.From(referenceInputAzureSql.Properties.FullSnapshotQuery))

				d.Set("delta_snapshot_query", pointer.From(referenceInputAzureSql.Properties.DeltaSnapshotQuery))

				d.Set("table", pointer.From(referenceInputAzureSql.Properties.Table))
			}
		}
	}

	return nil
}
