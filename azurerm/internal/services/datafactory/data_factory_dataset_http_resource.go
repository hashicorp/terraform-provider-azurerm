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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryDatasetHTTP() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryDatasetHTTPCreateUpdate,
		Read:   resourceArmDataFactoryDatasetHTTPRead,
		Update: resourceArmDataFactoryDatasetHTTPCreateUpdate,
		Delete: resourceArmDataFactoryDatasetHTTPDelete,

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

			// Specific to HTTP Dataset
			"relative_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Specific to HTTP Dataset
			"request_body": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Specific to HTTP Dataset
			"request_method": {
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
		},
	}
}

func resourceArmDataFactoryDatasetHTTPCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("Error checking for presence of existing Data Factory Dataset HTTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_dataset_http", *existing.ID)
		}
	}

	httpDatasetProperties := datafactory.HTTPDatasetTypeProperties{
		RelativeURL:   d.Get("relative_url").(string),
		RequestBody:   d.Get("request_body").(string),
		RequestMethod: d.Get("request_method").(string),
	}

	linkedServiceName := d.Get("linked_service_name").(string)
	linkedServiceType := "LinkedServiceReference"
	linkedService := &datafactory.LinkedServiceReference{
		ReferenceName: &linkedServiceName,
		Type:          &linkedServiceType,
	}

	description := d.Get("description").(string)
	// TODO
	httpTableset := datafactory.HTTPDataset{
		HTTPDatasetTypeProperties: &httpDatasetProperties,
		LinkedServiceName:         linkedService,
		Description:               &description,
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		httpTableset.Folder = &datafactory.DatasetFolder{
			Name: &name,
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		httpTableset.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		httpTableset.Annotations = &annotations
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		httpTableset.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("schema_column"); ok {
		httpTableset.Structure = expandDataFactoryDatasetStructure(v.([]interface{}))
	}

	datasetType := string(datafactory.TypeHTTPFile)
	dataset := datafactory.DatasetResource{
		Properties: &httpTableset,
		Type:       &datasetType,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, dataset, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Dataset HTTP  %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Dataset HTTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Dataset HTTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryDatasetHTTPRead(d, meta)
}

func resourceArmDataFactoryDatasetHTTPRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error retrieving Data Factory Dataset HTTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	httpTable, ok := resp.Properties.AsHTTPDataset()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Dataset HTTP %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeRelationalTable, *resp.Type)
	}

	d.Set("additional_properties", httpTable.AdditionalProperties)

	if httpTable.Description != nil {
		d.Set("description", httpTable.Description)
	}

	parameters := flattenDataFactoryParameters(httpTable.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	annotations := flattenDataFactoryAnnotations(httpTable.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	if linkedService := httpTable.LinkedServiceName; linkedService != nil {
		if linkedService.ReferenceName != nil {
			d.Set("linked_service_name", linkedService.ReferenceName)
		}
	}

	if properties := httpTable.HTTPDatasetTypeProperties; properties != nil {
		url, ok := properties.RelativeURL.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `relative_url` since it's not a string")
		} else {
			d.Set("relative_url", url)
		}
		requestBody, ok := properties.RequestBody.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `request_body` since it's not a string")
		} else {
			d.Set("request_body", requestBody)
		}
		requestMethod, ok := properties.RequestMethod.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping `request_method` since it's not a string")
		} else {
			d.Set("request_method", requestMethod)
		}
	}

	if folder := httpTable.Folder; folder != nil {
		if folder.Name != nil {
			d.Set("folder", folder.Name)
		}
	}

	structureColumns := flattenDataFactoryStructureColumns(httpTable.Structure)
	if err := d.Set("schema_column", structureColumns); err != nil {
		return fmt.Errorf("Error setting `schema_column`: %+v", err)
	}

	return nil
}

func resourceArmDataFactoryDatasetHTTPDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Data Factory Dataset HTTP %q (Data Factory %q / Resource Group %q): %s", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
