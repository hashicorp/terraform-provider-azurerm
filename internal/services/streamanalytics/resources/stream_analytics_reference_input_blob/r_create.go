// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_reference_input_blob

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsReferenceInputBlobCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := inputs.Input{
		Name: pointer.To(id.InputName),
		Properties: &inputs.ReferenceInputProperties{
			Datasource: &inputs.BlobReferenceInputDataSource{
				Properties: &inputs.BlobDataSourceProperties{
					Container:   pointer.To(d.Get("storage_container_name").(string)),
					DateFormat:  pointer.To(d.Get("date_format").(string)),
					PathPattern: pointer.To(d.Get("path_pattern").(string)),
					TimeFormat:  pointer.To(d.Get("time_format").(string)),
					StorageAccounts: &[]inputs.StorageAccount{
						{
							AccountName: pointer.To(d.Get("storage_account_name").(string)),
							AccountKey:  normalizeAccountKey(d.Get("storage_account_key").(string)),
						},
					},
					AuthenticationMode: pointer.ToEnum[inputs.AuthenticationMode](d.Get("authentication_mode").(string)),
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
