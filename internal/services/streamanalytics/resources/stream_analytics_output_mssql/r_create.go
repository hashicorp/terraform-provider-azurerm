// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_mssql

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputSqlCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).OutputsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_output_mssql", id.ID())
			}
		}
	}

	dataSourceProperties := outputs.AzureSqlDatabaseDataSourceProperties{
		Server:             pointer.To(d.Get("server").(string)),
		Database:           pointer.To(d.Get("database").(string)),
		Table:              pointer.To(d.Get("table").(string)),
		MaxBatchCount:      pointer.To(d.Get("max_batch_count").(float64)),
		MaxWriterCount:     pointer.To(d.Get("max_writer_count").(float64)),
		AuthenticationMode: pointer.ToEnum[outputs.AuthenticationMode](d.Get("authentication_mode").(string)),
	}

	// Add user/password dataSourceProperties only if authentication mode requires them
	if *dataSourceProperties.AuthenticationMode == outputs.AuthenticationModeConnectionString {
		dataSourceProperties.User = pointer.To(d.Get("user").(string))
		dataSourceProperties.Password = pointer.To(d.Get("password").(string))
	}

	props := outputs.Output{
		Name: pointer.To(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.AzureSqlDatabaseOutputDataSource{
				Properties: &dataSourceProperties,
			},
		},
	}

	var createOpts outputs.CreateOrReplaceOperationOptions
	var updateOpts outputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsOutputSqlRead(d, meta)
}
