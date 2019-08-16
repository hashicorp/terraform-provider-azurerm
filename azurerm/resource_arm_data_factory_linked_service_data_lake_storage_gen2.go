package azurerm

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryLinkedServiceDataLakeStorageGen2() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryLinkedServiceDataLakeStorageGen2CreateOrUpdate,
		Read:   resourceArmDataFactoryLinkedServiceDataLakeStorageGen2Read,
		Update: resourceArmDataFactoryLinkedServiceDataLakeStorageGen2CreateOrUpdate,
		Delete: resourceArmDataFactoryLinkedServiceDataLakeStorageGen2Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMDataFactoryLinkedServiceDatasetName,
			},

			"data_factory_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid name for Data Factory, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.URLIsHTTPS,
			},

			"service_principal_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.UUID,
			},

			"service_principal_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"tenant": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"integration_runtime_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"additional_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceArmDataFactoryLinkedServiceDataLakeStorageGen2CreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.LinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_data_lake_storage_gen2", *existing.ID)
		}
	}

	secureString := datafactory.SecureString{
		Value: utils.String(d.Get("service_principal_key").(string)),
		Type:  datafactory.TypeSecureString,
	}

	datalakeStorageGen2Properties := &datafactory.AzureBlobFSLinkedServiceTypeProperties{
		URL:                 utils.String(d.Get("url").(string)),
		ServicePrincipalID:  utils.String(d.Get("service_principal_id").(string)),
		Tenant:              utils.String(d.Get("tenant").(string)),
		ServicePrincipalKey: &secureString,
	}

	datalakeStorageGen2LinkedService := &datafactory.AzureBlobFSLinkedService{
		Description:                            utils.String(d.Get("description").(string)),
		AzureBlobFSLinkedServiceTypeProperties: datalakeStorageGen2Properties,
		Type:                                   datafactory.TypeAzureBlobFS,
	}

	if v, ok := d.GetOk("parameters"); ok {
		datalakeStorageGen2LinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		datalakeStorageGen2LinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		datalakeStorageGen2LinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		datalakeStorageGen2LinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: datalakeStorageGen2LinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDataFactoryLinkedServiceDataLakeStorageGen2Read(d, meta)
}

func resourceArmDataFactoryLinkedServiceDataLakeStorageGen2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.LinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["linkedservices"]

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	dataLakeStorageGen2, ok := resp.Properties.AsAzureBlobFSLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeAzureBlobFS, *resp.Type)
	}

	if dataLakeStorageGen2.Tenant != nil {
		d.Set("tenant", dataLakeStorageGen2.Tenant)
	}

	if dataLakeStorageGen2.ServicePrincipalID != nil {
		d.Set("service_principal_id", dataLakeStorageGen2.ServicePrincipalID)
	}

	if dataLakeStorageGen2.URL != nil {
		d.Set("url", dataLakeStorageGen2.URL)
	}

	d.Set("additional_properties", dataLakeStorageGen2.AdditionalProperties)

	if dataLakeStorageGen2.Description != nil {
		d.Set("description", dataLakeStorageGen2.Description)
	}

	annotations := flattenDataFactoryAnnotations(dataLakeStorageGen2.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(dataLakeStorageGen2.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := dataLakeStorageGen2.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceArmDataFactoryLinkedServiceDataLakeStorageGen2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactory.LinkedServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["linkedservices"]

	response, err := client.Delete(ctx, resourceGroup, dataFactoryName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("Error deleting Data Factory Linked Service Data Lake Storage Gen2 %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
