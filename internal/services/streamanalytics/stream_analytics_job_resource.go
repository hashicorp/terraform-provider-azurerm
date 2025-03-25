// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsJob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsJobCreate,
		Read:   resourceStreamAnalyticsJobRead,
		Update: resourceStreamAnalyticsJobUpdate,
		Delete: resourceStreamAnalyticsJobDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := streamingjobs.ParseStreamingJobID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsJobV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"stream_analytics_cluster_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"compatibility_level": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// NOTE: O+C Best Practice from MSFT is to use the latest version 1.2, but API uses 1.0 as the default
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					// values found in the other API the portal uses
					string(streamingjobs.CompatibilityLevelOnePointZero),
					"1.1",
					string(streamingjobs.CompatibilityLevelOnePointTwo),
				}, false),
			},

			"data_locale": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "en-US",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"events_late_arrival_max_delay_in_seconds": {
				// portal allows for up to 20d 23h 59m 59s
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(-1, 1814399),
				Default:      5,
			},

			"events_out_of_order_max_delay_in_seconds": {
				// portal allows for up to 9m 59s
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 599),
				Default:      0,
			},

			"events_out_of_order_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamingjobs.EventsOutOfOrderPolicyAdjust),
					string(streamingjobs.EventsOutOfOrderPolicyDrop),
				}, false),
				Default: string(streamingjobs.EventsOutOfOrderPolicyAdjust),
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamingjobs.JobTypeCloud),
					string(streamingjobs.JobTypeEdge),
				}, false),
				Default: string(streamingjobs.JobTypeCloud),
			},

			"output_error_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamingjobs.OutputErrorPolicyDrop),
					string(streamingjobs.OutputErrorPolicyStop),
				}, false),
				Default: string(streamingjobs.OutputErrorPolicyDrop),
			},

			"streaming_units": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 120),
			},

			"content_storage_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(streamingjobs.ContentStoragePolicySystemAccount),
				ValidateFunc: validation.StringInSlice([]string{
					string(streamingjobs.ContentStoragePolicySystemAccount),
					string(streamingjobs.ContentStoragePolicyJobStorageAccount),
				}, false),
			},

			"job_storage_account": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"authentication_mode": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(streamingjobs.AuthenticationModeConnectionString),
							ValidateFunc: validation.StringInSlice([]string{
								string(streamingjobs.AuthenticationModeConnectionString),
							}, false),
						},

						"account_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"account_key": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"transformation_query": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"job_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(streamingjobs.SkuNameStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(streamingjobs.SkuNameStandard),
					"StandardV2", // missing from swagger as described here https://github.com/Azure/azure-rest-api-specs/issues/27506
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceStreamAnalyticsJobCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Job creation.")

	id := streamingjobs.NewStreamingJobID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, id, streamingjobs.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_stream_analytics_job", id.ID())
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
			return fmt.Errorf("the job type `Edge` doesn't support `streaming_units`")
		}
	} else {
		if v, ok := d.GetOk("streaming_units"); ok {
			transformation.Properties.StreamingUnits = pointer.To(int64(v.(int)))
		} else {
			return fmt.Errorf("`streaming_units` must be set when `type` is `Cloud`")
		}
	}

	props := streamingjobs.StreamingJob{
		Name:     pointer.To(id.StreamingJobName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &streamingjobs.StreamingJobProperties{
			Sku: &streamingjobs.Sku{
				Name: pointer.To(streamingjobs.SkuName(d.Get("sku_name").(string))),
			},
			ContentStoragePolicy:               pointer.To(streamingjobs.ContentStoragePolicy(contentStoragePolicy)),
			EventsLateArrivalMaxDelayInSeconds: pointer.To(int64(d.Get("events_late_arrival_max_delay_in_seconds").(int))),
			EventsOutOfOrderMaxDelayInSeconds:  pointer.To(int64(d.Get("events_out_of_order_max_delay_in_seconds").(int))),
			EventsOutOfOrderPolicy:             pointer.To(streamingjobs.EventsOutOfOrderPolicy(d.Get("events_out_of_order_policy").(string))),
			OutputErrorPolicy:                  pointer.To(streamingjobs.OutputErrorPolicy(d.Get("output_error_policy").(string))),
			JobType:                            pointer.To(streamingjobs.JobType(jobType)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
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
		props.Properties.CompatibilityLevel = pointer.To(streamingjobs.CompatibilityLevel(d.Get("compatibility_level").(string)))
	}

	if v, ok := d.GetOk("job_storage_account"); ok {
		if contentStoragePolicy != string(streamingjobs.ContentStoragePolicyJobStorageAccount) {
			return fmt.Errorf("`job_storage_account` must not be set when `content_storage_policy` is `SystemAccount`")
		}
		props.Properties.JobStorageAccount = expandJobStorageAccount(v.([]interface{}))
	}

	if jobType == string(streamingjobs.JobTypeEdge) {
		if _, ok := d.GetOk("stream_analytics_cluster_id"); ok {
			return fmt.Errorf("the job type `Edge` doesn't support `stream_analytics_cluster_id`")
		}
	} else {
		if streamAnalyticsCluster := d.Get("stream_analytics_cluster_id"); streamAnalyticsCluster != "" {
			props.Properties.Cluster = &streamingjobs.ClusterInfo{
				Id: pointer.To(streamAnalyticsCluster.(string)),
			}
		} else {
			props.Properties.Cluster = &streamingjobs.ClusterInfo{
				Id: nil,
			}
		}
	}

	if dataLocale, ok := d.GetOk("data_locale"); ok {
		props.Properties.DataLocale = pointer.To(dataLocale.(string))
	}

	props.Properties.Transformation = &transformation

	if err := client.CreateOrReplaceThenPoll(ctx, id, props, streamingjobs.DefaultCreateOrReplaceOperationOptions()); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceStreamAnalyticsJobRead(d, meta)
}

func resourceStreamAnalyticsJobRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
				sku = string(pointer.From(props.Sku.Name))
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
			return tags.FlattenAndSet(d, model.Tags)
		}
	}
	return nil
}

func resourceStreamAnalyticsJobUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	transformationsClient := meta.(*clients.Client).StreamAnalytics.TransformationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Job update.")

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
				return fmt.Errorf("the job type `Edge` doesn't support `stream_analytics_cluster_id`")
			}
		}
		payload.Properties.Cluster = &streamingjobs.ClusterInfo{
			Id: pointer.To(clusterId),
		}
	}

	if d.HasChange("compatibility_level") {
		payload.Properties.CompatibilityLevel = pointer.To(streamingjobs.CompatibilityLevel(d.Get("compatibility_level").(string)))
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
		payload.Properties.EventsOutOfOrderPolicy = pointer.To(streamingjobs.EventsOutOfOrderPolicy(d.Get("events_out_of_order_policy").(string)))
	}

	if d.HasChange("output_error_policy") {
		payload.Properties.OutputErrorPolicy = pointer.To(streamingjobs.OutputErrorPolicy(d.Get("output_error_policy").(string)))
	}

	if d.HasChange("content_storage_policy") {
		payload.Properties.ContentStoragePolicy = pointer.To(streamingjobs.ContentStoragePolicy(d.Get("content_storage_policy").(string)))
	}

	if d.HasChange("job_storage_account") {
		storageAccount := d.Get("job_storage_account").([]interface{})
		if d.Get("content_storage_policy").(string) == string(streamingjobs.ContentStoragePolicyJobStorageAccount) {
			if len(storageAccount) == 0 {
				return fmt.Errorf("`job_storage_account` must be set when `content_storage_policy` is `JobStorageAccount`")
			}
		}
		payload.Properties.JobStorageAccount = expandJobStorageAccount(storageAccount)
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
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
			Name: pointer.To(streamingjobs.SkuName(d.Get("sku_name").(string))),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
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
				return fmt.Errorf("the job type `Edge` doesn't support `streaming_units`")
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

func resourceStreamAnalyticsJobDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := streamingjobs.ParseStreamingJobID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandJobStorageAccount(input []interface{}) *streamingjobs.JobStorageAccount {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	authenticationMode := v["authentication_mode"].(string)
	accountName := v["account_name"].(string)
	accountKey := v["account_key"].(string)

	return &streamingjobs.JobStorageAccount{
		AuthenticationMode: pointer.To(streamingjobs.AuthenticationMode(authenticationMode)),
		AccountName:        pointer.To(accountName),
		AccountKey:         pointer.To(accountKey),
	}
}

func flattenJobStorageAccount(d *pluginsdk.ResourceData, input *streamingjobs.JobStorageAccount) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	accountName := ""
	if v := input.AccountName; v != nil {
		accountName = *v
	}

	return []interface{}{
		map[string]interface{}{
			"authentication_mode": string(*input.AuthenticationMode),
			"account_name":        accountName,
			"account_key":         d.Get("job_storage_account.0.account_key").(string),
		},
	}
}
