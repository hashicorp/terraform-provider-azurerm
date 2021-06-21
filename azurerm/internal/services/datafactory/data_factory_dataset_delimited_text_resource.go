package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryDatasetDelimitedText() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetDelimitedTextCreateUpdate,
		Read:   resourceDataFactoryDatasetDelimitedTextRead,
		Update: resourceDataFactoryDatasetDelimitedTextCreateUpdate,
		Delete: resourceDataFactoryDatasetDelimitedTextDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"data_factory_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

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
						"filename": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
						"path": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"filename": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"row_delimiter": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"encoding": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"quote_character": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"escape_character": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"first_row_as_header": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			// Delimited Text Specific Field
			"null_value": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_delimited_text", *existing.ID)
		}
	}

	location := expandDataFactoryDatasetLocation(d)
	if location == nil {
		return fmt.Errorf("One of `http_server_location`, `azure_blob_storage_location` must be specified to create a DataFactory Delimited Text Dataset")
	}

	delimited_textDatasetProperties := datafactory.DelimitedTextDatasetTypeProperties{
		Location:         location,
		ColumnDelimiter:  d.Get("column_delimiter").(string),
		RowDelimiter:     d.Get("row_delimiter").(string),
		EncodingName:     d.Get("encoding").(string),
		QuoteChar:        d.Get("quote_character").(string),
		EscapeChar:       d.Get("escape_character").(string),
		FirstRowAsHeader: d.Get("first_row_as_header").(bool),
		NullValue:        d.Get("null_value").(string),
		CompressionLevel: d.Get("compression_level").(string),
		CompressionCodec: d.Get("compression_codec").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)
	// TODO
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
		delimited_textTableset.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, dataset, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Dataset DelimitedText  %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", id.FactoryName)

	delimited_textTable, ok := resp.Properties.AsDelimitedTextDataset()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", id.Name, id.FactoryName, id.ResourceGroup, datafactory.TypeBasicDatasetTypeDelimitedText, *resp.Type)
	}

	d.Set("additional_properties", delimited_textTable.AdditionalProperties)

	if delimited_textTable.Description != nil {
		d.Set("description", delimited_textTable.Description)
	}

	parameters := flattenDataFactoryParameters(delimited_textTable.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(delimited_textTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
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
		return fmt.Errorf("Error setting `schema_column`: %+v", err)
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
			return fmt.Errorf("Error deleting Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
		}
	}

	return nil
}
