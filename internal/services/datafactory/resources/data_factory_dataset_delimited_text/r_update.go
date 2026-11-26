// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_dataset_delimited_text

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

func resourceDataFactoryDatasetDelimitedTextCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).DataFactory.DatasetClient.SubscriptionID
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
				return tf.ImportAsExistsError("azurerm_data_factory_dataset_delimited_text", id.ID())
			}
		}
	}

	location := helpers.ExpandDataFactoryDatasetLocation(d)
	if location == nil {
		return fmt.Errorf("one of `http_server_location`, `azure_blob_fs_location`, `azure_blob_storage_location` must be specified to create a DataFactory Delimited Text Dataset")
	}

	delimited_textDatasetProperties := datafactory.DelimitedTextDatasetTypeProperties{
		Location: location,
	}

	if v, ok := d.Get("column_delimiter").(string); ok {
		delimited_textDatasetProperties.ColumnDelimiter = v
	}

	if v, ok := d.Get("row_delimiter").(string); ok {
		delimited_textDatasetProperties.RowDelimiter = v
	}

	if v, ok := d.Get("quote_character").(string); ok {
		delimited_textDatasetProperties.QuoteChar = v
	}

	if v, ok := d.Get("escape_character").(string); ok {
		delimited_textDatasetProperties.EscapeChar = v
	}

	if v, ok := d.GetOk("encoding"); ok {
		delimited_textDatasetProperties.EncodingName = v.(string)
	}

	if v, ok := d.Get("first_row_as_header").(bool); ok {
		delimited_textDatasetProperties.FirstRowAsHeader = v
	}

	if v, ok := d.Get("null_value").(string); ok {
		delimited_textDatasetProperties.NullValue = v
	}

	if v, ok := d.GetOk("compression_level"); ok {
		delimited_textDatasetProperties.CompressionLevel = v.(string)
	}

	if v, ok := d.GetOk("compression_codec"); ok {
		delimited_textDatasetProperties.CompressionCodec = v.(string)
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          pointer.To("LinkedServiceReference"),
	}

	description := d.Get("description").(string)
	delimited_textTableset := datafactory.DelimitedTextDataset{
		DelimitedTextDatasetTypeProperties: &delimited_textDatasetProperties,
		LinkedServiceName:                  linkedService,
		Description:                        &description,
	}

	if v, ok := d.GetOk("folder"); ok {
		delimited_textTableset.Folder = &datafactory.DatasetFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		delimited_textTableset.Parameters = helpers.ExpandDataSetParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("annotations"); ok {
		delimited_textTableset.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		delimited_textTableset.AdditionalProperties = v.(map[string]any)
	}

	if v, ok := d.GetOk("schema_column"); ok {
		delimited_textTableset.Structure = helpers.ExpandDataFactoryDatasetStructure(v.([]any))
	}

	dataset := datafactory.DatasetResource{
		Properties: &delimited_textTableset,
		Type:       pointer.To(string(datafactory.TypeBasicDatasetTypeDelimitedText)),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryDatasetDelimitedTextRead(d, meta)
}
