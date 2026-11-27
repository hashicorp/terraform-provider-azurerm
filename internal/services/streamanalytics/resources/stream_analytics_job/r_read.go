// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsJobRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingjobs.ParseStreamingJobID(d.Id())
	if err != nil {
		return err
	}

	opts := streamingjobs.GetOperationOptions{
		Expand: pointer.To("transformation"),
	}
	resp, err := client.Get(ctx, *id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("compatibility_level", pointer.From(props.CompatibilityLevel))
			d.Set("data_locale", pointer.From(props.DataLocale))
			d.Set("events_late_arrival_max_delay_in_seconds", pointer.From(props.EventsLateArrivalMaxDelayInSeconds))
			d.Set("events_out_of_order_max_delay_in_seconds", pointer.From(props.EventsOutOfOrderMaxDelayInSeconds))
			d.Set("events_out_of_order_policy", pointer.From(props.EventsOutOfOrderPolicy))
			d.Set("output_error_policy", pointer.From(props.OutputErrorPolicy))

			clusterId := ""
			if props.Cluster != nil && pointer.From(props.Cluster.Id) != "" {
				cId, err := clusters.ParseClusterID(*props.Cluster.Id)
				if err != nil {
					return err
				}
				clusterId = cId.ID()
			}
			d.Set("stream_analytics_cluster_id", clusterId)
			d.Set("type", pointer.From(props.JobType))

			sku := ""
			if props.Sku != nil {
				sku = pointer.FromEnum(props.Sku.Name)
			}
			d.Set("sku_name", sku)
			d.Set("content_storage_policy", pointer.From(props.ContentStoragePolicy))
			d.Set("job_id", pointer.From(props.JobId))
			d.Set("job_storage_account", flattenJobStorageAccount(d, props.JobStorageAccount))

			if transformation := props.Transformation; transformation != nil {
				if transformProps := transformation.Properties; transformProps != nil {
					d.Set("streaming_units", pointer.From(transformProps.StreamingUnits))
					d.Set("transformation_query", pointer.From(transformProps.Query))
				}
			}
			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}
	}
	return nil
}

func flattenJobStorageAccount(d *pluginsdk.ResourceData, input *streamingjobs.JobStorageAccount) []any {
	if input == nil {
		return []any{}
	}

	accountName := pointer.From(input.AccountName)

	return []any{
		map[string]any{
			"authentication_mode": string(*input.AuthenticationMode),
			"account_name":        accountName,
			"account_key":         d.Get("job_storage_account.0.account_key").(string),
		},
	}
}
