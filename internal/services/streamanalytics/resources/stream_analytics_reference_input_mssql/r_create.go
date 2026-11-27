// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_reference_input_mssql

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsReferenceInputMsSqlCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_reference_input_mssql", id.ID())
			}
		}
	}

	refreshType := d.Get("refresh_type").(string)

	if _, ok := d.GetOk("refresh_interval_duration"); refreshType != "Static" && !ok {
		return errors.New("refresh_interval_duration must be set if refresh_type is RefreshPeriodicallyWithFull or RefreshPeriodicallyWithDelta")
	} else if _, ok = d.GetOk("delta_snapshot_query"); refreshType == "Static" && ok {
		return errors.New("delta_snapshot_query cannot be set if refresh_type is Static")
	}

	properties := &inputs.AzureSqlReferenceInputDataSourceProperties{
		Server:      pointer.To(d.Get("server").(string)),
		Database:    pointer.To(d.Get("database").(string)),
		User:        pointer.To(d.Get("username").(string)),
		Password:    pointer.To(d.Get("password").(string)),
		RefreshType: pointer.ToEnum[inputs.RefreshType](refreshType),
	}

	if v, ok := d.GetOk("refresh_interval_duration"); ok {
		properties.RefreshRate = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("full_snapshot_query"); ok {
		properties.FullSnapshotQuery = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("delta_snapshot_query"); ok {
		properties.DeltaSnapshotQuery = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("table"); ok {
		properties.Table = pointer.To(v.(string))
	}

	props := inputs.Input{
		Name: pointer.To(id.InputName),
		Properties: &inputs.ReferenceInputProperties{
			Datasource: &inputs.AzureSqlReferenceInputDataSource{
				Properties: properties,
			},
		},
	}

	var opts inputs.CreateOrReplaceOperationOptions
	if _, err := client.CreateOrReplace(ctx, id, props, opts); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStreamAnalyticsReferenceInputMsSqlRead(d, meta)
}
