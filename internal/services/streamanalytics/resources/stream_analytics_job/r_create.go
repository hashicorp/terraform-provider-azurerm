// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsJobCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := streamingjobs.NewStreamingJobID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, streamingjobs.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_job", id.ID())
		}
	}

	// needs to be defined inline for a Create but via a separate API for Update
	transformation := streamingjobs.Transformation{
		Name: pointer.To("main"),
		Properties: &streamingjobs.TransformationProperties{
			Query: pointer.To(d.Get("transformation_query").(string)),
		},
	}

	contentStoragePolicy := d.Get("content_storage_policy").(string)
	jobType := d.Get("type").(string)

	if jobType == string(streamingjobs.JobTypeEdge) {
		if _, ok := d.GetOk("streaming_units"); ok {
			return errors.New("the job type `Edge` doesn't support `streaming_units`")
		}
	} else {
		if v, ok := d.GetOk("streaming_units"); ok {
			transformation.Properties.StreamingUnits = pointer.To(int64(v.(int)))
		} else {
			return errors.New("`streaming_units` must be set when `type` is `Cloud`")
		}
	}

	props := streamingjobs.StreamingJob{
		Name:     pointer.To(id.StreamingJobName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &streamingjobs.StreamingJobProperties{
			Sku: &streamingjobs.Sku{
				Name: pointer.ToEnum[streamingjobs.SkuName](d.Get("sku_name").(string)),
			},
			ContentStoragePolicy:               pointer.ToEnum[streamingjobs.ContentStoragePolicy](contentStoragePolicy),
			EventsLateArrivalMaxDelayInSeconds: pointer.To(int64(d.Get("events_late_arrival_max_delay_in_seconds").(int))),
			EventsOutOfOrderMaxDelayInSeconds:  pointer.To(int64(d.Get("events_out_of_order_max_delay_in_seconds").(int))),
			EventsOutOfOrderPolicy:             pointer.ToEnum[streamingjobs.EventsOutOfOrderPolicy](d.Get("events_out_of_order_policy").(string)),
			OutputErrorPolicy:                  pointer.ToEnum[streamingjobs.OutputErrorPolicy](d.Get("output_error_policy").(string)),
			JobType:                            pointer.ToEnum[streamingjobs.JobType](jobType),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	if expandedIdentity.Type == identity.TypeNone {
		// The StreamAnalytics API doesn't implement the standard `None` pattern - meaning that sending `None` outputs
		// an API error. This conditional is required to support this, else the API returns:
		//
		// >  Code="BadRequest" Message="The JSON provided in the request body is invalid. Cannot convert value 'None'
		// > to type 'System.Nullable`1[Microsoft.Streaming.Service.Contracts.CSMResourceProvider.IdentityType]"
		// > Details=[{"code":"400","correlationId":"dcdbdcfa-fe38-66f8-3aa3-36950bab0a28","message":"The JSON provided in the request body is invalid.
		// > Cannot convert value 'None' to type "System.Nullable`1[Microsoft.Streaming.Service.Contracts.CSMResourceProvider.IdentityType]"
		//
		// Tracked in https://github.com/Azure/azure-rest-api-specs/issues/17649
		expandedIdentity = nil
	}
	props.Identity = expandedIdentity

	if _, ok := d.GetOk("compatibility_level"); ok {
		props.Properties.CompatibilityLevel = pointer.ToEnum[streamingjobs.CompatibilityLevel](d.Get("compatibility_level").(string))
	}

	if v, ok := d.GetOk("job_storage_account"); ok {
		props.Properties.JobStorageAccount = expandJobStorageAccount(v.([]any))
	}

	if jobType == string(streamingjobs.JobTypeEdge) {
		if _, ok := d.GetOk("stream_analytics_cluster_id"); ok {
			return errors.New("the job type `Edge` doesn't support `stream_analytics_cluster_id`")
		}
	} else {
		if streamAnalyticsCluster := d.Get("stream_analytics_cluster_id"); streamAnalyticsCluster != "" {
			props.Properties.Cluster = &streamingjobs.ClusterInfo{
				Id: pointer.To(streamAnalyticsCluster.(string)),
			}
		} else {
			props.Properties.Cluster = &streamingjobs.ClusterInfo{}
		}
	}

	if dataLocale, ok := d.GetOk("data_locale"); ok {
		props.Properties.DataLocale = pointer.To(dataLocale.(string))
	}

	props.Properties.Transformation = &transformation

	if err := client.CreateOrReplaceCallbackThenPoll(ctx, id, props, streamingjobs.DefaultCreateOrReplaceOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceStreamAnalyticsJobRead(d, meta)
}

func expandJobStorageAccount(input []any) *streamingjobs.JobStorageAccount {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	jobStorageAccount := streamingjobs.JobStorageAccount{
		AuthenticationMode: pointer.ToEnum[streamingjobs.AuthenticationMode](v["authentication_mode"].(string)),
		AccountName:        pointer.To(v["account_name"].(string)),
	}

	if accountKey := v["account_key"].(string); accountKey != "" {
		jobStorageAccount.AccountKey = pointer.To(accountKey)
	}

	return &jobStorageAccount
}
