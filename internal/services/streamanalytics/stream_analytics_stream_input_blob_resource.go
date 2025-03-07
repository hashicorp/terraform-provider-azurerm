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
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsStreamInputBlob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsStreamInputBlobCreateUpdate,
		Read:   resourceStreamAnalyticsStreamInputBlobRead,
		Update: resourceStreamAnalyticsStreamInputBlobCreateUpdate,
		Delete: resourceStreamAnalyticsStreamInputBlobDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := inputs.ParseInputID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsStreamInputBlobV0ToV1{},
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

			"storage_account_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
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

			"serialization": schemaStreamAnalyticsStreamInputSerialization(),

			"authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(inputs.AuthenticationModeConnectionString),
				ValidateFunc: validation.StringInSlice([]string{
					string(inputs.AuthenticationModeConnectionString),
					string(inputs.AuthenticationModeMsi),
				}, false),
			},
		},
	}
}

func resourceStreamAnalyticsStreamInputBlobCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Stream Input Blob creation.")
	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_stream_input_blob", id.ID())
		}
	}

	containerName := d.Get("storage_container_name").(string)
	dateFormat := d.Get("date_format").(string)
	pathPattern := d.Get("path_pattern").(string)
	storageAccountKey := d.Get("storage_account_key").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	timeFormat := d.Get("time_format").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := inputs.Input{
		Name: utils.String(id.InputName),
		Properties: &inputs.StreamInputProperties{
			// Type: streamanalytics.TypeBasicInputPropertiesTypeStream,
			Datasource: &inputs.BlobStreamInputDataSource{
				// Type: streamanalytics.TypeBasicStreamInputDataSourceTypeMicrosoftStorageBlob,
				Properties: &inputs.BlobStreamInputDataSourceProperties{
					Container:   utils.String(containerName),
					DateFormat:  utils.String(dateFormat),
					PathPattern: utils.String(pathPattern),
					TimeFormat:  utils.String(timeFormat),
					StorageAccounts: &[]inputs.StorageAccount{
						{
							AccountName: utils.String(storageAccountName),
							AccountKey:  utils.String(storageAccountKey),
						},
					},
					AuthenticationMode: pointer.To(inputs.AuthenticationMode(d.Get("authentication_mode").(string))),
				},
			},
			Serialization: serialization,
		},
	}

	var createOpts inputs.CreateOrReplaceOperationOptions
	var updateOpts inputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsStreamInputBlobRead(d, meta)
}

func resourceStreamAnalyticsStreamInputBlobRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			input, ok := props.(inputs.InputProperties) // nolint: gosimple
			if !ok {
				return fmt.Errorf("converting %s to an Input", *id)
			}

			streamInput, ok := input.(inputs.StreamInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to a Stream Input", *id)
			}

			streamBlobInput, ok := streamInput.Datasource.(inputs.BlobStreamInputDataSource)
			if !ok {
				return fmt.Errorf("converting Stream Input Blob to an Stream Input: %+v", err)
			}

			if streamBlobInputProps := streamBlobInput.Properties; streamBlobInputProps != nil {
				dateFormat := ""
				if v := streamBlobInput.Properties.DateFormat; v != nil {
					dateFormat = *v
				}
				d.Set("date_format", dateFormat)

				pathPattern := ""
				if v := streamBlobInputProps.PathPattern; v != nil {
					pathPattern = *v
				}
				d.Set("path_pattern", pathPattern)

				containerName := ""
				if v := streamBlobInputProps.Container; v != nil {
					containerName = *v
				}
				d.Set("storage_container_name", containerName)

				timeFormat := ""
				if v := streamBlobInputProps.TimeFormat; v != nil {
					timeFormat = *v
				}
				d.Set("time_format", timeFormat)

				authMode := ""
				if v := streamBlobInputProps.AuthenticationMode; v != nil {
					authMode = string(*v)
				}
				d.Set("authentication_mode", authMode)

				if accounts := streamBlobInputProps.StorageAccounts; accounts != nil && len(*accounts) > 0 {
					account := (*accounts)[0]
					d.Set("storage_account_name", account.AccountName)
				}

				if err := d.Set("serialization", flattenStreamAnalyticsStreamInputSerialization(streamInput.Serialization)); err != nil {
					return fmt.Errorf("setting `serialization`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceStreamAnalyticsStreamInputBlobDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
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
