// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_mssql

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputSqlRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OutputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			output, ok := props.Datasource.(outputs.AzureSqlDatabaseOutputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a SQL Output", *id)
			}

			d.Set("server", pointer.From(output.Properties.Server))

			d.Set("database", pointer.From(output.Properties.Database))

			d.Set("table", pointer.From(output.Properties.Table))

			d.Set("user", pointer.From(output.Properties.User))

			authMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authMode = string(*v)
			}
			d.Set("authentication_mode", authMode)

			maxBatchCount := float64(10000)
			if v := output.Properties.MaxBatchCount; v != nil {
				maxBatchCount = *v
			}
			d.Set("max_batch_count", maxBatchCount)

			maxWriterCount := float64(1)
			if v := output.Properties.MaxWriterCount; v != nil {
				maxWriterCount = *v
			}
			d.Set("max_writer_count", maxWriterCount)
		}
	}
	return nil
}
