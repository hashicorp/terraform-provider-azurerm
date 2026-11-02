// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_binary

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryDatasetBinaryRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	binaryTable, ok := resp.Properties.AsBinaryDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset Binary %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeBinary, resp.Properties)
	}

	d.Set("additional_properties", binaryTable.AdditionalProperties)

	if binaryTable.Description != nil {
		d.Set("description", binaryTable.Description)
	}

	if err := d.Set("parameters", helpers.FlattenDataSetParameters(binaryTable.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(binaryTable.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := binaryTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := binaryTable.BinaryDatasetTypeProperties; properties != nil {
		if httpServerLocation, ok := properties.Location.AsHTTPServerLocation(); ok {
			if err := d.Set("http_server_location", helpers.FlattenDataFactoryDatasetHTTPServerLocation(httpServerLocation)); err != nil {
				return fmt.Errorf("setting `http_server_location` for Data Factory Binary Dataset %s", err)
			}
		}
		if azureBlobStorageLocation, ok := properties.Location.AsAzureBlobStorageLocation(); ok {
			if err := d.Set("azure_blob_storage_location", helpers.FlattenDataFactoryDatasetAzureBlobStorageLocation(azureBlobStorageLocation)); err != nil {
				return fmt.Errorf("setting `azure_blob_storage_location` for Data Factory Binary Dataset %s", err)
			}
		}
		if sftpLocation, ok := properties.Location.AsSftpLocation(); ok {
			if err := d.Set("sftp_server_location", helpers.FlattenDataFactoryDatasetSFTPLocation(sftpLocation)); err != nil {
				return fmt.Errorf("setting `sftp_server_location` for Data Factory Binary Dataset %s", err)
			}
		}

		if err := d.Set("compression", helpers.FlattenDataFactoryDatasetCompression(properties.Compression)); err != nil {
			return fmt.Errorf("setting `compression`: %+v", err)
		}
	}

	if folder := binaryTable.Folder; folder != nil && folder.Name != nil {
		d.Set("folder", folder.Name)
	}

	return nil
}
