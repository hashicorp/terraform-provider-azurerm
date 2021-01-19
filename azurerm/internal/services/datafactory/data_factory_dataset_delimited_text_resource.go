package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryDatasetDelimitedText() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataFactoryDatasetDelimitedTextCreateUpdate,
		Read:   resourceDataFactoryDatasetDelimitedTextRead,
		Update: resourceDataFactoryDatasetDelimitedTextCreateUpdate,
		Delete: resourceDataFactoryDatasetDelimitedTextDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMDataFactoryLinkedServiceDatasetName,
			},

			"data_factory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"linked_service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field, one option for 'location'
			"http_server_location": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConflictsWith: []string{"sftp_server_location", "file_server_location", "s3_location", "azure_blob_storage_location"},
				ConflictsWith: []string{"azure_blob_storage_location"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relative_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"filename": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			// Delimited Text Specific Field, one option for 'location'
			"azure_blob_storage_location": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConflictsWith: []string{"sftp_server_location", "file_server_location", "s3_location", "azure_blob_storage_location"},
				ConflictsWith: []string{"http_server_location"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"filename": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			// Delimited Text Specific Field
			"column_delimiter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"row_delimiter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"encoding": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"quote_character": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"escape_character": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Delimited Text Specific Field
			"first_row_as_header": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Delimited Text Specific Field
			"null_value": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"folder": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"additional_properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"schema_column": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:     schema.TypeString,
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
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"compression_level": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Optimal",
					"Fastest",
				}, false),
			},
		},
	}
}

func resourceDataFactoryDatasetDelimitedTextCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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

	datasetType := string(datafactory.TypeDelimitedText)
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

func resourceDataFactoryDatasetDelimitedTextRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error classifiying Data Factory Dataset DelimitedText %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", id.Name, id.FactoryName, id.ResourceGroup, datafactory.TypeRelationalTable, *resp.Type)
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
		if httpServerLocation, ok := properties.Location.AsHTTPServerLocation(); ok {
			if err := d.Set("http_server_location", flattenDataFactoryDatasetHTTPServerLocation(httpServerLocation)); err != nil {
				return fmt.Errorf("Error setting `http_server_location` for Data Factory Delimited Text Dataset %s", err)
			}
		}
		if azureBlobStorageLocation, ok := properties.Location.AsAzureBlobStorageLocation(); ok {
			if err := d.Set("azure_blob_storage_location", flattenDataFactoryDatasetAzureBlobStorageLocation(azureBlobStorageLocation)); err != nil {
				return fmt.Errorf("Error setting `azure_blob_storage_location` for Data Factory Delimited Text Dataset %s", err)
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

func resourceDataFactoryDatasetDelimitedTextDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandDataFactoryDatasetLocation(d *schema.ResourceData) datafactory.BasicDatasetLocation {
	if _, ok := d.GetOk("http_server_location"); ok {
		return expandDataFactoryDatasetHttpServerLocation(d)
	}

	if _, ok := d.GetOk("azure_blob_storage_location"); ok {
		return expandDataFactoryDatasetAzureBlobStorageLocation(d)
	}

	return nil
}

func expandDataFactoryDatasetHttpServerLocation(d *schema.ResourceData) datafactory.BasicDatasetLocation {
	props := d.Get("http_server_location").([]interface{})[0].(map[string]interface{})
	relativeUrl := props["relative_url"].(string)
	path := props["path"].(string)
	filename := props["filename"].(string)

	httpServerLocation := datafactory.HTTPServerLocation{
		RelativeURL: relativeUrl,
		FolderPath:  path,
		FileName:    filename,
	}
	return httpServerLocation
}

func expandDataFactoryDatasetAzureBlobStorageLocation(d *schema.ResourceData) datafactory.BasicDatasetLocation {
	props := d.Get("azure_blob_storage_location").([]interface{})[0].(map[string]interface{})
	container := props["container"].(string)
	path := props["path"].(string)
	filename := props["filename"].(string)

	blobStorageLocation := datafactory.AzureBlobStorageLocation{
		Container:  container,
		FolderPath: path,
		FileName:   filename,
	}
	return blobStorageLocation
}

func flattenDataFactoryDatasetHTTPServerLocation(input *datafactory.HTTPServerLocation) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.RelativeURL != nil {
		result["relative_url"] = input.RelativeURL
	}
	if input.FolderPath != nil {
		result["path"] = input.FolderPath
	}
	if input.FileName != nil {
		result["filename"] = input.FileName
	}

	return []interface{}{result}
}

func flattenDataFactoryDatasetAzureBlobStorageLocation(input *datafactory.AzureBlobStorageLocation) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if input.Container != nil {
		result["container"] = input.Container
	}
	if input.FolderPath != nil {
		result["path"] = input.FolderPath
	}
	if input.FileName != nil {
		result["filename"] = input.FileName
	}

	return []interface{}{result}
}
