// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

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

func resourceStreamAnalyticsReferenceInputBlob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsReferenceInputBlobCreate,
		Read:   resourceStreamAnalyticsReferenceInputBlobRead,
		Update: resourceStreamAnalyticsReferenceInputBlobUpdate,
		Delete: resourceStreamAnalyticsReferenceInputBlobDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := inputs.ParseInputID(id)
			return err
		}, importStreamAnalyticsReferenceInput("Microsoft.Storage/Blob")),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsReferenceInputBlobV0ToV1{},
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
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

func resourceStreamAnalyticsReferenceInputBlobCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Reference Input Blob creation.")
	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_reference_input_blob", id.ID())
		}
	}

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := inputs.Input{
		Name: utils.String(id.InputName),
		Properties: &inputs.ReferenceInputProperties{
			Datasource: &inputs.BlobReferenceInputDataSource{
				Properties: &inputs.BlobDataSourceProperties{
					Container:   utils.String(d.Get("storage_container_name").(string)),
					DateFormat:  utils.String(d.Get("date_format").(string)),
					PathPattern: utils.String(d.Get("path_pattern").(string)),
					TimeFormat:  utils.String(d.Get("time_format").(string)),
					StorageAccounts: &[]inputs.StorageAccount{
						{
							AccountName: utils.String(d.Get("storage_account_name").(string)),
							AccountKey:  normalizeAccountKey(d.Get("storage_account_key").(string)),
						},
					},
					AuthenticationMode: utils.ToPtr(inputs.AuthenticationMode(d.Get("authentication_mode").(string))),
				},
			},
			Serialization: serialization,
		},
	}

	var opts inputs.CreateOrReplaceOperationOptions
	if _, err := client.CreateOrReplace(ctx, id, props, opts); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStreamAnalyticsReferenceInputBlobRead(d, meta)
}

func resourceStreamAnalyticsReferenceInputBlobUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Reference Input Blob update.")
	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	// TODO d.HasChanges()
	props := inputs.Input{
		Name: utils.String(id.InputName),
		Properties: &inputs.ReferenceInputProperties{
			Datasource: &inputs.BlobReferenceInputDataSource{
				Properties: &inputs.BlobDataSourceProperties{
					Container:   utils.String(d.Get("storage_container_name").(string)),
					DateFormat:  utils.String(d.Get("date_format").(string)),
					PathPattern: utils.String(d.Get("path_pattern").(string)),
					TimeFormat:  utils.String(d.Get("time_format").(string)),
					StorageAccounts: &[]inputs.StorageAccount{
						{
							AccountName: utils.String(d.Get("storage_account_name").(string)),
							AccountKey:  normalizeAccountKey(d.Get("storage_account_key").(string)),
						},
					},
					AuthenticationMode: utils.ToPtr(inputs.AuthenticationMode(d.Get("authentication_mode").(string))),
				},
			},
			Serialization: serialization,
		},
	}

	var opts inputs.UpdateOperationOptions
	if _, err := client.Update(ctx, *id, props, opts); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceStreamAnalyticsReferenceInputBlobRead(d, meta)
}

func resourceStreamAnalyticsReferenceInputBlobRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
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

			dataSource, ok := input.(inputs.ReferenceInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to a Reference Input", *id)
			}

			referenceInputBlob, ok := dataSource.Datasource.(inputs.BlobReferenceInputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a Blob Reference Input", *id)
			}

			if referenceInputBlob.Properties != nil {
				dateFormat := ""
				if v := referenceInputBlob.Properties.DateFormat; v != nil {
					dateFormat = *v
				}
				d.Set("date_format", dateFormat)

				pathPattern := ""
				if v := referenceInputBlob.Properties.PathPattern; v != nil {
					pathPattern = *v
				}
				d.Set("path_pattern", pathPattern)

				containerName := ""
				if v := referenceInputBlob.Properties.Container; v != nil {
					containerName = *v
				}
				d.Set("storage_container_name", containerName)

				timeFormat := ""
				if v := referenceInputBlob.Properties.TimeFormat; v != nil {
					timeFormat = *v
				}
				d.Set("time_format", timeFormat)

				authMode := ""
				if v := referenceInputBlob.Properties.AuthenticationMode; v != nil {
					authMode = string(*v)
				}
				d.Set("authentication_mode", authMode)

				if accounts := referenceInputBlob.Properties.StorageAccounts; accounts != nil && len(*accounts) > 0 {
					account := (*accounts)[0]
					d.Set("storage_account_name", account.AccountName)
				}
			}
			if err := d.Set("serialization", flattenStreamAnalyticsStreamInputSerialization(dataSource.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
		}
	}

	return nil
}

func resourceStreamAnalyticsReferenceInputBlobDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func normalizeAccountKey(accountKey string) *string {
	if accountKey != "" {
		return utils.String(accountKey)
	}

	return nil
}
