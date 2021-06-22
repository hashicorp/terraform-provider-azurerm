package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryDatasetSFTP() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetSFTPCreateUpdate,
		Read:   resourceDataFactoryDatasetSFTPRead,
		Update: resourceDataFactoryDatasetSFTPCreateUpdate,
		Delete: resourceDataFactoryDatasetSFTPDelete,

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

			// SFTP Specific field
			"path": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// SFTP Specific field
			"filename": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"compression": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						//TarGZip, GZip, ZipDeflate
						"level": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						// SFTP Specific field
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(datafactory.TypeBasicDatasetCompressionTypeBZip2),
								string(datafactory.TypeBasicDatasetCompressionTypeDeflate),
								string(datafactory.TypeBasicDatasetCompressionTypeGZip),
								string(datafactory.TypeBasicDatasetCompressionTypeTar),
								string(datafactory.TypeBasicDatasetCompressionTypeTarGZip),
								string(datafactory.TypeBasicDatasetCompressionTypeZipDeflate),
							}, false),
						},
					},
				},
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
		},
	}
}

func resourceDataFactoryDatasetSFTPCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for presence of existing Data Factory Dataset SFTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_sftp", *existing.ID)
		}
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)

	fileShareDataset := datafactory.FileShareDataset{
		FileShareDatasetTypeProperties: &datafactory.FileShareDatasetTypeProperties{
			FolderPath: d.Get("path").(string),
			FileName:   d.Get("filename").(string),
		},
		LinkedServiceName: linkedService,
		Description:       &description,
	}

	if _, ok := d.GetOk("compression"); ok {
		fileShareDataset.Compression = expandDataFactoryCompression(d)
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		fileShareDataset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		fileShareDataset.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		fileShareDataset.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		fileShareDataset.AdditionalProperties = v.(map[string]interface{})
	}

	datasetType := string(datafactory.TypeBasicDatasetTypeFileShare)
	dataset := datafactory.DatasetResource{
		Properties: &fileShareDataset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Dataset SFTP  %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Data Factory Dataset SFTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Data Factory Dataset SFTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryDatasetSFTPRead(d, meta)
}

func resourceDataFactoryDatasetSFTPRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["datasets"]

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Data Factory Dataset SFTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	fileShareDataSet, ok := resp.Properties.AsFileShareDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset SFTP %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeBasicDatasetTypeFileShare, *resp.Type)
	}

	d.Set("additional_properties", fileShareDataSet.AdditionalProperties)

	if fileShareDataSet.Description != nil {
		d.Set("description", fileShareDataSet.Description)
	}

	parameters := flattenDataFactoryParameters(fileShareDataSet.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(fileShareDataSet.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := fileShareDataSet.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := fileShareDataSet.FileShareDatasetTypeProperties; properties != nil {
		filename, ok := properties.FileName.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `filename` since it's not a string")
		} else {
			d.Set("filename", filename)
		}
		path, ok := properties.FolderPath.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `path` since it's not a string")
		} else {
			d.Set("path", path)
		}
		compression := flattenDataFactoryCompression(fileShareDataSet.Compression)
		if err := d.Set("compression", compression); err != nil {
			return fmt.Errorf("setting `compression`: %+v", err)
		}

	}

	if folder := fileShareDataSet.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	return nil
}

func resourceDataFactoryDatasetSFTPDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["datasets"]

	response, err := client.Delete(ctx, resourceGroup, dataFactoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Data Factory Dataset SFTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}

func flattenDataFactoryCompression(input datafactory.BasicDatasetCompression) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	if compression, ok := input.AsDatasetBZip2Compression(); ok {
		result["type"] = compression.Type
	}
	if compression, ok := input.AsDatasetDeflateCompression(); ok {
		result["type"] = compression.Type
	}
	if compression, ok := input.AsDatasetGZipCompression(); ok {
		result["type"] = compression.Type
		result["level"] = compression.Level
	}
	if compression, ok := input.AsDatasetTarCompression(); ok {
		result["type"] = compression.Type
	}
	if compression, ok := input.AsDatasetTarGZipCompression(); ok {
		result["type"] = compression.Type
		result["level"] = compression.Level
	}
	if compression, ok := input.AsDatasetZipDeflateCompression(); ok {
		result["type"] = compression.Type
		result["level"] = compression.Level
	}

	return []interface{}{result}
}

func expandDataFactoryCompression(d *pluginsdk.ResourceData) datafactory.BasicDatasetCompression {
	props := d.Get("compression").([]interface{})[0].(map[string]interface{})
	level := props["level"].(string)
	compressionType := props["type"].(string)

	if datafactory.TypeBasicDatasetCompression(compressionType) == datafactory.TypeBasicDatasetCompressionTypeBZip2 {
		return datafactory.DatasetBZip2Compression{
			Type: datafactory.TypeBasicDatasetCompression(compressionType),
		}
	}
	if datafactory.TypeBasicDatasetCompression(compressionType) == datafactory.TypeBasicDatasetCompressionTypeDeflate {
		return datafactory.DatasetDeflateCompression{
			Type: datafactory.TypeBasicDatasetCompression(compressionType),
		}
	}
	if datafactory.TypeBasicDatasetCompression(compressionType) == datafactory.TypeBasicDatasetCompressionTypeGZip {
		return datafactory.DatasetGZipCompression{
			Type:  datafactory.TypeBasicDatasetCompression(compressionType),
			Level: level,
		}
	}
	if datafactory.TypeBasicDatasetCompression(compressionType) == datafactory.TypeBasicDatasetCompressionTypeTar {
		return datafactory.DatasetTarCompression{
			Type: datafactory.TypeBasicDatasetCompression(compressionType),
		}
	}
	if datafactory.TypeBasicDatasetCompression(compressionType) == datafactory.TypeBasicDatasetCompressionTypeTarGZip {
		return datafactory.DatasetTarGZipCompression{
			Type:  datafactory.TypeBasicDatasetCompression(compressionType),
			Level: level,
		}
	}
	if datafactory.TypeBasicDatasetCompression(compressionType) == datafactory.TypeBasicDatasetCompressionTypeZipDeflate {
		return datafactory.DatasetZipDeflateCompression{
			Type:  datafactory.TypeBasicDatasetCompression(compressionType),
			Level: level,
		}
	}

	return nil
}
