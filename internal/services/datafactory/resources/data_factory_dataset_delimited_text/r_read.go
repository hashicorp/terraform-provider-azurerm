// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_delimited_text

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryDatasetDelimitedTextRead(d *pluginsdk.ResourceData, meta any) error {
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

	d.Set("name", resp.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	delimited_textTable, ok := resp.Properties.AsDelimitedTextDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset DelimitedText %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeDelimitedText, resp.Properties)
	}

	d.Set("additional_properties", delimited_textTable.AdditionalProperties)

	if delimited_textTable.Description != nil {
		d.Set("description", delimited_textTable.Description)
	}

	if err := d.Set("parameters", helpers.FlattenDataSetParameters(delimited_textTable.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(delimited_textTable.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := delimited_textTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := delimited_textTable.DelimitedTextDatasetTypeProperties; properties != nil {
		switch location := properties.Location.(type) {
		case datafactory.HTTPServerLocation:
			if err := d.Set("http_server_location", helpers.FlattenDataFactoryDatasetHTTPServerLocation(&location)); err != nil {
				return fmt.Errorf("setting `http_server_location` for Data Factory Delimited Text Dataset %s", err)
			}
		case datafactory.AzureBlobStorageLocation:
			if err := d.Set("azure_blob_storage_location", helpers.FlattenDataFactoryDatasetAzureBlobStorageLocation(&location)); err != nil {
				return fmt.Errorf("setting `azure_blob_storage_location` for Data Factory Delimited Text Dataset %s", err)
			}
		case datafactory.AzureBlobFSLocation:
			if err := d.Set("azure_blob_fs_location", helpers.FlattenDataFactoryDatasetAzureBlobFSLocation(&location)); err != nil {
				return fmt.Errorf("setting `azure_blob_fs_location` for Data Factory Delimited Text Dataset %s", err)
			}
		}

		columnDelimiter, ok := properties.ColumnDelimiter.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `column_delimiter` since it's not a string")
		} else {
			d.Set("column_delimiter", columnDelimiter)
		}

		rowDelimiter, ok := properties.RowDelimiter.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `row_delimiter` since it's not a string")
		} else {
			d.Set("row_delimiter", rowDelimiter)
		}

		encodingName, ok := properties.EncodingName.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `encoding` since it's not a string")
		} else {
			d.Set("encoding", encodingName)
		}

		quoteChar, ok := properties.QuoteChar.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `quote_char` since it's not a string")
		} else {
			d.Set("quote_character", quoteChar)
		}

		escapeChar, ok := properties.EscapeChar.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `escape_char` since it's not a string")
		} else {
			d.Set("escape_character", escapeChar)
		}
		firstRow, ok := properties.FirstRowAsHeader.(bool)
		if !ok {
			log.Printf("[DEBUG] Skipping `first_row_as_header` since it's not a string")
		} else {
			d.Set("first_row_as_header", firstRow)
		}
		nullValue, ok := properties.NullValue.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `null_value` since it's not a string")
		} else {
			d.Set("null_value", nullValue)
		}
		compressionLevel, ok := properties.CompressionLevel.(string)
		if !ok {
			log.Printf("[DEBUG] skipping `compression_level` since it's not a string")
		} else {
			d.Set("compression_level", compressionLevel)
		}
		compressionCodec, ok := properties.CompressionCodec.(string)
		if !ok {
			log.Printf("[DEBUG] skipping `compression_codec` since it's not a string")
		} else {
			d.Set("compression_codec", compressionCodec)
		}
	}

	if folder := delimited_textTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	if err := d.Set("schema_column", helpers.FlattenDataFactoryStructureColumns(delimited_textTable.Structure)); err != nil {
		return fmt.Errorf("setting `schema_column`: %+v", err)
	}

	return nil
}
