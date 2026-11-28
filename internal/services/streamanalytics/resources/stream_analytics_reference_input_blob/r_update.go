// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_reference_input_blob

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsReferenceInputBlobUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	// TODO d.HasChanges()
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

	var opts inputs.UpdateOperationOptions
	if _, err := client.Update(ctx, *id, props, opts); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceStreamAnalyticsReferenceInputBlobRead(d, meta)
}
