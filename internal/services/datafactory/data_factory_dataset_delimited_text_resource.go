// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryDatasetDelimitedText() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetDelimitedTextCreateUpdate,
		Read:   resourceDataFactoryDatasetDelimitedTextRead,
		Update: resourceDataFactoryDatasetDelimitedTextCreateUpdate,
		Delete: resourceDataFactoryDatasetDelimitedTextDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
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
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"linked_service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field, one option for 'location'
			"http_server_location": {
				Type:         pluginsdk.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"http_server_location", "azure_blob_storage_location", "azure_blob_fs_location"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"relative_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"path": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"dynamic_path_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"filename": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"dynamic_filename_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			// Delimited Text Specific Field, one option for 'location'
			"azure_blob_storage_location": {
				Type:         pluginsdk.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"http_server_location", "azure_blob_storage_location", "azure_blob_fs_location"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"container": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"dynamic_container_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"dynamic_path_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"filename": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"dynamic_filename_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"azure_blob_fs_location": {
				Type:         pluginsdk.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"http_server_location", "azure_blob_storage_location", "azure_blob_fs_location"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"file_system": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"filename": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			// Delimited Text Specific Field
			"column_delimiter": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  ",",
			},

			// Delimited Text Specific Field
			"row_delimiter": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			// Delimited Text Specific Field
			"quote_character": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  `"`,
			},

			// Delimited Text Specific Field
			"escape_character": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  `\`,
			},

			// Delimited Text Specific Field
			"encoding": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"first_row_as_header": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Delimited Text Specific Field
			"null_value": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "",
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"schema_column": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Byte",
								"Byte[]",
								"Boolean",
								"Date",
								"DateTime",
								"DateTimeOffset",
								"Decimal",
								"Double",
								"Guid",
								"Int16",
								"Int32",
								"Int64",
								"Single",
								"String",
								"TimeSpan",
							}, false),
						},
						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"compression_codec": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"bzip2",
					"gzip",
					"deflate",
					"ZipDeflate",
					"TarGzip",
					"Tar",
					"snappy",
					"lz4",
				}, false),
			},

			"compression_level": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Optimal",
					"Fastest",
				}, false),
			},
		},
	}
}

func resourceDataFactoryDatasetDelimitedTextCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_delimited_text", id.ID())
		}
	}

	location := expandDataFactoryDatasetLocation(d)
	if location == nil {
		return fmt.Errorf("one of `http_server_location`, `azure_blob_storage_location` must be specified to create a DataFactory Delimited Text Dataset")
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
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)
	delimited_textTableset := datafactory.DelimitedTextDataset{
		DelimitedTextDatasetTypeProperties: &delimited_textDatasetProperties,
		LinkedServiceName:                  linkedService,
		Description:                        &description,
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		delimited_textTableset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		delimited_textTableset.Parameters = expandDataSetParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		delimited_textTableset.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		delimited_textTableset.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("schema_column"); ok {
		delimited_textTableset.Structure = expandDataFactoryDatasetStructure(v.([]interface{}))
	}

	datasetType := string(datafactory.TypeBasicDatasetTypeDelimitedText)
	dataset := datafactory.DatasetResource{
		Properties: &delimited_textTableset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryDatasetDelimitedTextRead(d, meta)
}

func resourceDataFactoryDatasetDelimitedTextRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		if utils.ResponseWasNotFound(resp.Response) {
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

	parameters := flattenDataSetParameters(delimited_textTable.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(delimited_textTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
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
			if err := d.Set("http_server_location", flattenDataFactoryDatasetHTTPServerLocation(&location)); err != nil {
				return fmt.Errorf("setting `http_server_location` for Data Factory Delimited Text Dataset %s", err)
			}
		case datafactory.AzureBlobStorageLocation:
			if err := d.Set("azure_blob_storage_location", flattenDataFactoryDatasetAzureBlobStorageLocation(&location)); err != nil {
				return fmt.Errorf("setting `azure_blob_storage_location` for Data Factory Delimited Text Dataset %s", err)
			}
		case datafactory.AzureBlobFSLocation:
			if err := d.Set("azure_blob_fs_location", flattenDataFactoryDatasetAzureBlobFSLocation(&location)); err != nil {
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

	structureColumns := flattenDataFactoryStructureColumns(delimited_textTable.Structure)
	if err := d.Set("schema_column", structureColumns); err != nil {
		return fmt.Errorf("setting `schema_column`: %+v", err)
	}

	return nil
}

func resourceDataFactoryDatasetDelimitedTextDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
