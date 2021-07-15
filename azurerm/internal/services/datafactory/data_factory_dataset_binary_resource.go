package datafactory

import (
	"fmt"
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

func resourceDataFactoryDatasetBinary() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetBinaryCreateUpdate,
		Read:   resourceDataFactoryDatasetBinaryRead,
		Update: resourceDataFactoryDatasetBinaryCreateUpdate,
		Delete: resourceDataFactoryDatasetBinaryDelete,
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

			// TODO: replace with `data_factory_id` in 3.0
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

			// Binary Dataset Specific Field
			"http_server_location": {
				Type:          pluginsdk.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"azure_blob_storage_location", "sftp_server_location"},
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

			"sftp_server_location": {
				Type:          pluginsdk.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"azure_blob_storage_location", "http_server_location"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"path": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"filename": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			// Binary Dataset Specific Field
			"azure_blob_storage_location": {
				Type:          pluginsdk.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"http_server_location", "sftp_server_location"},
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

			"compression": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// TarGZip, GZip, ZipDeflate
						"level": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Optimal",
								"Fastest",
							}, false),
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
		},
	}
}

func resourceDataFactoryDatasetBinaryCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDataSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("data_factory_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory Dataset Binary %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_binary", *existing.ID)
		}
	}

	location := expandDataFactoryDatasetLocation(d)
	if location == nil {
		return fmt.Errorf("one of `http_server_location`, `azure_blob_storage_location` or `sftp_server_location`, must be specified to create a DataFactory Binary Dataset")
	}

	binaryDatasetProperties := datafactory.BinaryDatasetTypeProperties{
		Location: location,
	}

	if _, ok := d.GetOk("compression"); ok {
		binaryDatasetProperties.Compression = expandDataFactoryDatasetCompression(d)
	}

	binaryTableset := datafactory.BinaryDataset{
		BinaryDatasetTypeProperties: &binaryDatasetProperties,
		Description:                 utils.String(d.Get("description").(string)),
		LinkedServiceName: &datafactory.LinkedServiceReference{
			ReferenceName: utils.String(d.Get("linked_service_name").(string)),
			Type:          utils.String("LinkedServiceReference"),
		},
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		binaryTableset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		binaryTableset.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		binaryTableset.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		binaryTableset.AdditionalProperties = v.(map[string]interface{})
	}

	datasetType := string(datafactory.TypeBasicDatasetTypeBinary)
	dataset := datafactory.DatasetResource{
		Properties: &binaryTableset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating Data Factory Dataset Binary  %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	if _, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, ""); err != nil {
		return fmt.Errorf("retrieving Data Factory Dataset Binary %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryDatasetBinaryRead(d, meta)
}

func resourceDataFactoryDatasetBinaryRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("retrieving Data Factory Dataset Binary %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", id.FactoryName)

	binaryTable, ok := resp.Properties.AsBinaryDataset()
	if !ok {
		return fmt.Errorf("classifiying Data Factory Dataset Binary %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", id.Name, id.FactoryName, id.ResourceGroup, datafactory.TypeBasicDatasetTypeBinary, *resp.Type)
	}

	d.Set("additional_properties", binaryTable.AdditionalProperties)

	if binaryTable.Description != nil {
		d.Set("description", binaryTable.Description)
	}

	if err := d.Set("parameters", flattenDataFactoryParameters(binaryTable.Parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(binaryTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := binaryTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := binaryTable.BinaryDatasetTypeProperties; properties != nil {
		if httpServerLocation, ok := properties.Location.AsHTTPServerLocation(); ok {
			if err := d.Set("http_server_location", flattenDataFactoryDatasetHTTPServerLocation(httpServerLocation)); err != nil {
				return fmt.Errorf("setting `http_server_location` for Data Factory Binary Dataset %s", err)
			}
		}
		if azureBlobStorageLocation, ok := properties.Location.AsAzureBlobStorageLocation(); ok {
			if err := d.Set("azure_blob_storage_location", flattenDataFactoryDatasetAzureBlobStorageLocation(azureBlobStorageLocation)); err != nil {
				return fmt.Errorf("setting `azure_blob_storage_location` for Data Factory Binary Dataset %s", err)
			}
		}
		if sftpLocation, ok := properties.Location.AsSftpLocation(); ok {
			if err := d.Set("sftp_server_location", flattenDataFactoryDatasetSFTPLocation(sftpLocation)); err != nil {
				return fmt.Errorf("setting `sftp_server_location` for Data Factory Binary Dataset %s", err)
			}
		}

		compression := flattenDataFactoryDatasetCompression(properties.Compression)
		if err := d.Set("compression", compression); err != nil {
			return fmt.Errorf("setting `compression`: %+v", err)
		}
	}

	if folder := binaryTable.Folder; folder != nil && folder.Name != nil {
		d.Set("folder", folder.Name)
	}

	return nil
}

func resourceDataFactoryDatasetBinaryDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("deleting Data Factory Dataset Binary %q (Data Factory %q / Resource Group %q): %s", id.Name, id.FactoryName, id.ResourceGroup, err)
		}
	}

	return nil
}
