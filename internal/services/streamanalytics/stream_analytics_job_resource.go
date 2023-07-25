// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/transformations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsJob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsJobCreateUpdate,
		Read:   resourceStreamAnalyticsJobRead,
		Update: resourceStreamAnalyticsJobCreateUpdate,
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
				Computed:     true,
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
				ValidateFunc: validate.StreamAnalyticsJobStreamingUnits,
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

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"job_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceStreamAnalyticsJobCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	transformationsClient := meta.(*clients.Client).StreamAnalytics.TransformationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Job creation.")

	id := streamingjobs.NewStreamingJobID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if d.IsNewResource() {
		var opts streamingjobs.GetOperationOptions
		existing, err := client.Get(ctx, id, opts)
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
		Name: utils.String("main"),
		Properties: &streamingjobs.TransformationProperties{
			Query: utils.String(d.Get("transformation_query").(string)),
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
			transformation.Properties.StreamingUnits = utils.Int64(int64(v.(int)))
		} else {
			return fmt.Errorf("`streaming_units` must be set when `type` is `Cloud`")
		}
	}

	expandedIdentity, err := expandStreamAnalyticsJobIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props := streamingjobs.StreamingJob{
		Name:     utils.String(id.StreamingJobName),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &streamingjobs.StreamingJobProperties{
			Sku: &streamingjobs.Sku{
				Name: utils.ToPtr(streamingjobs.SkuNameStandard),
			},
			ContentStoragePolicy:               utils.ToPtr(streamingjobs.ContentStoragePolicy(contentStoragePolicy)),
			EventsLateArrivalMaxDelayInSeconds: utils.Int64(int64(d.Get("events_late_arrival_max_delay_in_seconds").(int))),
			EventsOutOfOrderMaxDelayInSeconds:  utils.Int64(int64(d.Get("events_out_of_order_max_delay_in_seconds").(int))),
			EventsOutOfOrderPolicy:             utils.ToPtr(streamingjobs.EventsOutOfOrderPolicy(d.Get("events_out_of_order_policy").(string))),
			OutputErrorPolicy:                  utils.ToPtr(streamingjobs.OutputErrorPolicy(d.Get("output_error_policy").(string))),
			JobType:                            utils.ToPtr(streamingjobs.JobType(jobType)),
		},
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("compatibility_level"); ok {
		compatibilityLevel := d.Get("compatibility_level").(string)
		props.Properties.CompatibilityLevel = utils.ToPtr(streamingjobs.CompatibilityLevel(compatibilityLevel))
	}

	if contentStoragePolicy == string(streamingjobs.ContentStoragePolicyJobStorageAccount) {
		if v, ok := d.GetOk("job_storage_account"); ok {
			props.Properties.JobStorageAccount = expandJobStorageAccount(v.([]interface{}))
		} else {
			return fmt.Errorf("`job_storage_account` must be set when `content_storage_policy` is `JobStorageAccount`")
		}
	}

	if jobType == string(streamingjobs.JobTypeEdge) {
		if _, ok := d.GetOk("stream_analytics_cluster_id"); ok {
			return fmt.Errorf("the job type `Edge` doesn't support `stream_analytics_cluster_id`")
		}
	} else {
		if streamAnalyticsCluster := d.Get("stream_analytics_cluster_id"); streamAnalyticsCluster != "" {
			props.Properties.Cluster = &streamingjobs.ClusterInfo{
				Id: utils.String(streamAnalyticsCluster.(string)),
			}
		} else {
			props.Properties.Cluster = &streamingjobs.ClusterInfo{
				Id: nil,
			}
		}
	}

	if dataLocale, ok := d.GetOk("data_locale"); ok {
		props.Properties.DataLocale = utils.String(dataLocale.(string))
	}

	if d.IsNewResource() {
		props.Properties.Transformation = &transformation

		var opts streamingjobs.CreateOrReplaceOperationOptions
		if err := client.CreateOrReplaceThenPoll(ctx, id, props, opts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else {
		var updateOpts streamingjobs.UpdateOperationOptions
		if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		if d.HasChanges("streaming_units", "transformation_query") {
			transformationUpdate := transformations.Transformation{
				Name: utils.String("main"),
				Properties: &transformations.TransformationProperties{
					Query: utils.String(d.Get("transformation_query").(string)),
				},
			}

			if jobType == string(streamingjobs.JobTypeEdge) {
				if _, ok := d.GetOk("streaming_units"); ok {
					return fmt.Errorf("the job type `Edge` doesn't support `streaming_units`")
				}
			} else {
				if v, ok := d.GetOk("streaming_units"); ok {
					transformationUpdate.Properties.StreamingUnits = utils.Int64(int64(v.(int)))
				} else {
					return fmt.Errorf("`streaming_units` must be set when `type` is `Cloud`")
				}
			}

			transformationId := transformations.NewTransformationID(subscriptionId, id.ResourceGroupName, id.StreamingJobName, *transformation.Name)

			var updateOpts transformations.UpdateOperationOptions
			if _, err := transformationsClient.Update(ctx, transformationId, transformationUpdate, updateOpts); err != nil {
				return fmt.Errorf("updating transformation for %s: %+v", id, err)
			}
		}
	}

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
		Expand: utils.ToPtr("transformation"),
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

		if err := d.Set("identity", flattenJobIdentity(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %v", err)
		}
		if props := model.Properties; props != nil {
			compatibilityLevel := ""
			if v := props.CompatibilityLevel; v != nil {
				compatibilityLevel = string(*v)
			}
			d.Set("compatibility_level", compatibilityLevel)

			dataLocale := ""
			if v := props.DataLocale; v != nil {
				dataLocale = *v
			}
			d.Set("data_locale", dataLocale)

			var lateArrival int64
			if v := props.EventsLateArrivalMaxDelayInSeconds; v != nil {
				lateArrival = *v
			}
			d.Set("events_late_arrival_max_delay_in_seconds", lateArrival)

			var maxDelay int64
			if v := props.EventsOutOfOrderMaxDelayInSeconds; v != nil {
				maxDelay = *v
			}
			d.Set("events_out_of_order_max_delay_in_seconds", maxDelay)

			orderPolicy := ""
			if v := props.EventsOutOfOrderPolicy; v != nil {
				orderPolicy = string(*v)
			}
			d.Set("events_out_of_order_policy", orderPolicy)

			outputPolicy := ""
			if v := props.OutputErrorPolicy; v != nil {
				outputPolicy = string(*v)
			}
			d.Set("output_error_policy", outputPolicy)

			cluster := ""
			if props.Cluster != nil && props.Cluster.Id != nil {
				cluster = *props.Cluster.Id
			}
			d.Set("stream_analytics_cluster_id", cluster)

			jobType := ""
			if v := props.JobType; v != nil {
				jobType = string(*v)
			}
			d.Set("type", jobType)

			storagePolicy := ""
			if v := props.ContentStoragePolicy; v != nil {
				storagePolicy = string(*v)
			}
			d.Set("content_storage_policy", storagePolicy)

			jobId := ""
			if v := props.JobId; v != nil {
				jobId = *v
			}
			d.Set("job_id", jobId)

			d.Set("job_storage_account", flattenJobStorageAccount(d, props.JobStorageAccount))

			if transformation := props.Transformation; transformation != nil {
				var streamingUnits int64
				if v := props.Transformation.Properties.StreamingUnits; v != nil {
					streamingUnits = *v
				}
				d.Set("streaming_units", streamingUnits)

				query := ""
				if v := props.Transformation.Properties.Query; v != nil {
					query = *v
				}
				d.Set("transformation_query", query)
			}
			return tags.FlattenAndSet(d, model.Tags)
		}
	}
	return nil
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

func expandStreamAnalyticsJobIdentity(input []interface{}) (*streamingjobs.Identity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	// Otherwise we get:
	//   Code="BadRequest"
	//   Message="The JSON provided in the request body is invalid. Cannot convert value 'None' to
	//   type 'System.Nullable`1[Microsoft.Streaming.Service.Contracts.CSMResourceProvider.IdentityType]"
	// Upstream issue: https://github.com/Azure/azure-rest-api-specs/issues/17649
	if expanded.Type == identity.TypeNone {
		return nil, nil
	}

	return &streamingjobs.Identity{
		Type: utils.String(string(expanded.Type)),
	}, nil
}

func flattenJobIdentity(identity *streamingjobs.Identity) []interface{} {
	if identity == nil {
		return nil
	}

	var t string
	if identity.Type != nil {
		t = *identity.Type
	}

	var tenantId string
	if identity.TenantId != nil {
		tenantId = *identity.TenantId
	}

	var principalId string
	if identity.PrincipalId != nil {
		principalId = *identity.PrincipalId
	}

	return []interface{}{
		map[string]interface{}{
			"type":         t,
			"tenant_id":    tenantId,
			"principal_id": principalId,
		},
	}
}

func expandJobStorageAccount(input []interface{}) *streamingjobs.JobStorageAccount {
	if input == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	authenticationMode := v["authentication_mode"].(string)
	accountName := v["account_name"].(string)
	accountKey := v["account_key"].(string)

	return &streamingjobs.JobStorageAccount{
		AuthenticationMode: utils.ToPtr(streamingjobs.AuthenticationMode(authenticationMode)),
		AccountName:        utils.String(accountName),
		AccountKey:         utils.String(accountKey),
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
