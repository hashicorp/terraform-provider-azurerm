// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryDatasetJSON() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryDatasetJSONCreateUpdate,
		Read:   resourceDataFactoryDatasetJSONRead,
		Update: resourceDataFactoryDatasetJSONCreateUpdate,
		Delete: resourceDataFactoryDatasetJSONDelete,

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

			// JSON Dataset Specific Field
			"http_server_location": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConflictsWith: []string{"sftp_server_location", "file_server_location", "s3_location", "azure_blob_storage_location"},
				ConflictsWith: []string{"azure_blob_storage_location"},
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

			// JSON Dataset Specific Field
			"azure_blob_storage_location": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConflictsWith: []string{"sftp_server_location", "file_server_location", "s3_location", "azure_blob_storage_location"},
				ConflictsWith: []string{"http_server_location"},
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
							Required: true,
						},
						"dynamic_path_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"filename": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"dynamic_filename_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			// JSON Dataset Specific Field
			"encoding": {
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
		},
	}
}

func resourceDataFactoryDatasetJSONCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	jsonDatasetProperties := datafactory.JSONDatasetTypeProperties{
		Location:     location,
		EncodingName: d.Get("encoding").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)
	jsonTableset := datafactory.JSONDataset{
		JSONDatasetTypeProperties: &jsonDatasetProperties,
		LinkedServiceName:         linkedService,
		Description:               &description,
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		jsonTableset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		jsonTableset.Parameters = expandDataSetParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		jsonTableset.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		jsonTableset.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("schema_column"); ok {
		jsonTableset.Structure = expandDataFactoryDatasetStructure(v.([]interface{}))
	}

	datasetType := string(datafactory.TypeBasicDatasetTypeJSON)
	dataset := datafactory.DatasetResource{
		Properties: &jsonTableset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryDatasetJSONRead(d, meta)
}

func resourceDataFactoryDatasetJSONRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

	d.Set("name", id.Name)
	d.Set("data_factory_id", dataFactoryId.ID())

	jsonTable, ok := resp.Properties.AsJSONDataset()
	if !ok {
		return fmt.Errorf("classifying Data Factory Dataset JSON %s: Expected: %q Received: %T", *id, datafactory.TypeBasicDatasetTypeJSON, resp.Properties)
	}

	d.Set("additional_properties", jsonTable.AdditionalProperties)

	if jsonTable.Description != nil {
		d.Set("description", jsonTable.Description)
	}

	parameters := flattenDataSetParameters(jsonTable.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(jsonTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	if linkedService := jsonTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := jsonTable.JSONDatasetTypeProperties; properties != nil {
		if httpServerLocation, ok := properties.Location.AsHTTPServerLocation(); ok {
			if err := d.Set("http_server_location", flattenDataFactoryDatasetHTTPServerLocation(httpServerLocation)); err != nil {
				return fmt.Errorf("setting `http_server_location` for Data Factory Delimited Text Dataset %s", err)
			}
		}
		if azureBlobStorageLocation, ok := properties.Location.AsAzureBlobStorageLocation(); ok {
			if err := d.Set("azure_blob_storage_location", flattenDataFactoryDatasetAzureBlobStorageLocation(azureBlobStorageLocation)); err != nil {
				return fmt.Errorf("setting `azure_blob_storage_location` for Data Factory Delimited Text Dataset %s", err)
			}
		}

		encodingName, ok := properties.EncodingName.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `encoding` since it's not a string")
		} else {
			d.Set("encoding", encodingName)
		}
	}

	if folder := jsonTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	structureColumns := flattenDataFactoryStructureColumns(jsonTable.Structure)
	if err := d.Set("schema_column", structureColumns); err != nil {
		return fmt.Errorf("setting `schema_column`: %+v", err)
	}

	return nil
}

func resourceDataFactoryDatasetJSONDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
