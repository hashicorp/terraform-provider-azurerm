// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsJobUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).JobsClient
	transformationsClient := streamanalyticsClient.From(meta).TransformationsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := streamingjobs.ParseStreamingJobID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, *id, streamingjobs.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", err)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", err)
	}

	payload := existing.Model

	if d.HasChange("stream_analytics_cluster_id") {
		clusterId := d.Get("stream_analytics_cluster_id").(string)
		if d.Get("type").(string) == string(streamingjobs.JobTypeEdge) {
			if clusterId != "" {
				return errors.New("the job type `Edge` doesn't support `stream_analytics_cluster_id`")
			}
		}
		payload.Properties.Cluster = &streamingjobs.ClusterInfo{
			Id: pointer.To(clusterId),
		}
	}

	if d.HasChange("compatibility_level") {
		payload.Properties.CompatibilityLevel = pointer.ToEnum[streamingjobs.CompatibilityLevel](d.Get("compatibility_level").(string))
	}

	if d.HasChange("data_locale") {
		payload.Properties.DataLocale = pointer.To(d.Get("data_locale").(string))
	}

	if d.HasChange("events_late_arrival_max_delay_in_seconds") {
		payload.Properties.EventsLateArrivalMaxDelayInSeconds = pointer.To(int64(d.Get("events_late_arrival_max_delay_in_seconds").(int)))
	}

	if d.HasChange("events_out_of_order_max_delay_in_seconds") {
		payload.Properties.EventsOutOfOrderMaxDelayInSeconds = pointer.To(int64(d.Get("events_out_of_order_max_delay_in_seconds").(int)))
	}

	if d.HasChange("events_out_of_order_policy") {
		payload.Properties.EventsOutOfOrderPolicy = pointer.ToEnum[streamingjobs.EventsOutOfOrderPolicy](d.Get("events_out_of_order_policy").(string))
	}

	if d.HasChange("output_error_policy") {
		payload.Properties.OutputErrorPolicy = pointer.ToEnum[streamingjobs.OutputErrorPolicy](d.Get("output_error_policy").(string))
	}

	if d.HasChange("content_storage_policy") {
		payload.Properties.ContentStoragePolicy = pointer.ToEnum[streamingjobs.ContentStoragePolicy](d.Get("content_storage_policy").(string))
	}

	if d.HasChange("job_storage_account") {
		storageAccount := d.Get("job_storage_account").([]any)
		if d.Get("content_storage_policy").(string) == string(streamingjobs.ContentStoragePolicyJobStorageAccount) {
			if len(storageAccount) == 0 {
				return errors.New("`job_storage_account` must be set when `content_storage_policy` is `JobStorageAccount`")
			}
		}
		payload.Properties.JobStorageAccount = expandJobStorageAccount(storageAccount)
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		if expandedIdentity.Type == identity.TypeNone {
			// See comment in create, tracked in https://github.com/Azure/azure-rest-api-specs/issues/17649
			expandedIdentity = nil
		}
		payload.Identity = expandedIdentity
	}

	if d.HasChange("sku_name") {
		payload.Properties.Sku = &streamingjobs.Sku{
			Name: pointer.ToEnum[streamingjobs.SkuName](d.Get("sku_name").(string)),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.Update(ctx, *id, *payload, streamingjobs.DefaultUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.HasChanges("transformation_query", "streaming_units") {
		transformation := transformations.Transformation{
			Name: pointer.To("main"),
			Properties: &transformations.TransformationProperties{
				Query: pointer.To(d.Get("transformation_query").(string)),
			},
		}

		streamingUnits := d.Get("streaming_units").(int)
		if d.Get("type").(string) == string(streamingjobs.JobTypeEdge) {
			if streamingUnits != 0 {
				return errors.New("the job type `Edge` doesn't support `streaming_units`")
			}
		}
		transformation.Properties.StreamingUnits = pointer.To(int64(streamingUnits))

		transformationId := transformations.NewTransformationID(subscriptionId, id.ResourceGroupName, id.StreamingJobName, *transformation.Name)

		if _, err := transformationsClient.Update(ctx, transformationId, transformation, transformations.DefaultUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating transformation for %s: %+v", id, err)
		}
	}

	return resourceStreamAnalyticsJobRead(d, meta)
}
