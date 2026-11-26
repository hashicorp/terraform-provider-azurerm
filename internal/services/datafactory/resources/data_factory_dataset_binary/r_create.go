// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_binary

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryDatasetBinaryCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_dataset_binary", id.ID())
			}
		}
	}

	location := helpers.ExpandDataFactoryDatasetLocation(d)
	if location == nil {
		return fmt.Errorf("one of `http_server_location`, `azure_blob_storage_location` or `sftp_server_location`, must be specified to create a DataFactory Binary Dataset")
	}

	binaryDatasetProperties := datafactory.BinaryDatasetTypeProperties{
		Location: location,
	}

	if _, ok := d.GetOk("compression"); ok {
		binaryDatasetProperties.Compression = helpers.ExpandDataFactoryDatasetCompression(d)
	}

	binaryTableset := datafactory.BinaryDataset{
		BinaryDatasetTypeProperties: &binaryDatasetProperties,
		Description:                 pointer.To(d.Get("description").(string)),
		LinkedServiceName: &datafactory.LinkedServiceReference{
			ReferenceName: pointer.To(d.Get("linked_service_name").(string)),
			Type:          pointer.To("LinkedServiceReference"),
		},
	}

	if v, ok := d.GetOk("folder"); ok {
		binaryTableset.Folder = &datafactory.DatasetFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		binaryTableset.Parameters = helpers.ExpandDataSetParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("annotations"); ok {
		binaryTableset.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		binaryTableset.AdditionalProperties = v.(map[string]any)
	}

	dataset := datafactory.DatasetResource{
		Properties: &binaryTableset,
		Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeBinary)),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryDatasetBinaryRead(d, meta)
}
