// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsOutputBlob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsOutputBlobCreateUpdate,
		Read:   resourceStreamAnalyticsOutputBlobRead,
		Update: resourceStreamAnalyticsOutputBlobCreateUpdate,
		Delete: resourceStreamAnalyticsOutputBlobDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := outputs.ParseOutputID(id)
			return err
		}, importStreamAnalyticsOutput(outputs.BlobOutputDataSource{})),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputBlobV0ToV1{},
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

			"stream_analytics_job_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"date_format": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"path_pattern": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"time_format": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"serialization": schemaStreamAnalyticsOutputSerialization(),

			"authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(outputs.AuthenticationModeConnectionString),
				ValidateFunc: validation.StringInSlice([]string{
					string(outputs.AuthenticationModeConnectionString),
					string(outputs.AuthenticationModeMsi),
				}, false),
			},

			"batch_max_wait_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.BatchMaxWaitTime,
			},
			"batch_min_rows": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1000000),
			},

			"storage_account_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceStreamAnalyticsOutputBlobCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_output_blob", id.ID())
		}
	}

	containerName := d.Get("storage_container_name").(string)
	dateFormat := d.Get("date_format").(string)
	pathPattern := d.Get("path_pattern").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	timeFormat := d.Get("time_format").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := outputs.Output{
		Name: utils.String(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.BlobOutputDataSource{
				Properties: &outputs.BlobOutputDataSourceProperties{
					StorageAccounts: &[]outputs.StorageAccount{
						{
							AccountKey:  getStorageAccountKey(d.Get("storage_account_key").(string)),
							AccountName: utils.String(storageAccountName),
						},
					},
					Container:          utils.String(containerName),
					DateFormat:         utils.String(dateFormat),
					PathPattern:        utils.String(pathPattern),
					TimeFormat:         utils.String(timeFormat),
					AuthenticationMode: utils.ToPtr(outputs.AuthenticationMode(d.Get("authentication_mode").(string))),
				},
			},
			Serialization: serialization,
		},
	}

	if batchMaxWaitTime, ok := d.GetOk("batch_max_wait_time"); ok {
		props.Properties.TimeWindow = utils.String(batchMaxWaitTime.(string))
	}

	if batchMinRows, ok := d.GetOk("batch_min_rows"); ok {
		props.Properties.SizeWindow = utils.Int64(int64(batchMinRows.(int)))
	}

	// timeWindow and sizeWindow must be set for Parquet serialization
	_, isParquet := serialization.(outputs.ParquetSerialization)
	if isParquet && (props.Properties.TimeWindow == nil || props.Properties.SizeWindow == nil) {
		return fmt.Errorf("cannot create %s: batch_min_rows and batch_time_window must be set for Parquet serialization", id)
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

	return resourceStreamAnalyticsOutputBlobRead(d, meta)
}

func resourceStreamAnalyticsOutputBlobRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			output, ok := props.Datasource.(outputs.BlobOutputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a Blob Output", *id)
			}

			dateFormat := ""
			if v := output.Properties.DateFormat; v != nil {
				dateFormat = *v
			}
			d.Set("date_format", dateFormat)

			pathPattern := ""
			if v := output.Properties.PathPattern; v != nil {
				pathPattern = *v
			}
			d.Set("path_pattern", pathPattern)

			containerName := ""
			if v := output.Properties.Container; v != nil {
				containerName = *v
			}
			d.Set("storage_container_name", containerName)

			timeFormat := ""
			if v := output.Properties.TimeFormat; v != nil {
				timeFormat = *v
			}
			d.Set("time_format", timeFormat)

			authenticationMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authenticationMode = string(*v)
			}
			d.Set("authentication_mode", authenticationMode)

			if accounts := output.Properties.StorageAccounts; accounts != nil && len(*accounts) > 0 {
				account := (*accounts)[0]
				d.Set("storage_account_name", account.AccountName)
			}

			if err := d.Set("serialization", flattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
			d.Set("batch_max_wait_time", props.TimeWindow)
			d.Set("batch_min_rows", props.SizeWindow)
		}
	}
	return nil
}

func resourceStreamAnalyticsOutputBlobDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func getStorageAccountKey(input string) *string {
	if input == "" {
		return nil
	}

	return utils.String(input)
}
