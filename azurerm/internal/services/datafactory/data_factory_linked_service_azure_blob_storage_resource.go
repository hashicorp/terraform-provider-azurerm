package datafactory

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryLinkedServiceAzureBlobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataFactoryLinkedServiceBlobStorageCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceBlobStorageRead,
		Update: resourceDataFactoryLinkedServiceBlobStorageCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceBlobStorageDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validate.LinkedServiceDatasetName,
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

			"connection_string": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"connection_string", "sas_uri", "service_endpoint"},
			},

			"sas_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"connection_string", "sas_uri", "service_endpoint"},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"integration_runtime_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"use_managed_identity": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ConflictsWith: []string{
					"service_principal_id",
				},
			},

			"service_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"use_managed_identity"},
				ExactlyOneOf: []string{"connection_string", "sas_uri", "service_endpoint"},
			},

			"service_principal_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"service_principal_key"},
				ConflictsWith: []string{
					"use_managed_identity",
				},
			},

			"service_principal_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				RequiredWith: []string{"service_principal_id"},
			},

			"tenant_id": {
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDataFactoryLinkedServiceBlobStorageCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service BlobStorage Anonymous %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_azure_blob_storage", *existing.ID)
		}
	}

	blobStorageProperties := &datafactory.AzureBlobStorageLinkedServiceTypeProperties{}

	if v, ok := d.GetOk("connection_string"); ok {
		blobStorageProperties.ConnectionString = &datafactory.SecureString{
			Value: utils.String(v.(string)),
			Type:  datafactory.TypeTypeSecureString,
		}
	}

	if v, ok := d.GetOk("sas_uri"); ok {
		blobStorageProperties.SasURI = &datafactory.SecureString{
			Value: utils.String(v.(string)),
			Type:  datafactory.TypeTypeSecureString,
		}
	}

	if d.Get("use_managed_identity").(bool) {
		if v, ok := d.GetOk("service_endpoint"); ok {
			blobStorageProperties.ServiceEndpoint = utils.String(v.(string))
		}
	} else {
		secureString := datafactory.SecureString{
			Value: utils.String(d.Get("service_principal_key").(string)),
			Type:  datafactory.TypeTypeSecureString,
		}

		blobStorageProperties.ServicePrincipalID = utils.String(d.Get("service_principal_id").(string))
		blobStorageProperties.Tenant = utils.String(d.Get("tenant_id").(string))
		blobStorageProperties.ServicePrincipalKey = &secureString
	}

	blobStorageLinkedService := &datafactory.AzureBlobStorageLinkedService{
		Description: utils.String(d.Get("description").(string)),
		AzureBlobStorageLinkedServiceTypeProperties: blobStorageProperties,
		Type: datafactory.TypeBasicLinkedServiceTypeAzureBlobStorage,
	}

	if v, ok := d.GetOk("parameters"); ok {
		blobStorageLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		blobStorageLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		blobStorageLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		blobStorageLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: blobStorageLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service BlobStorage %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service BlobStorage %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service BlobStorage %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryLinkedServiceBlobStorageRead(d, meta)
}

func resourceDataFactoryLinkedServiceBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

		return fmt.Errorf("Error retrieving Data Factory Linked Service BlobStorage %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	blobStorage, ok := resp.Properties.AsAzureBlobStorageLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service BlobStorage %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeBasicLinkedServiceTypeAzureBlobStorage, *resp.Type)
	}

	if blobStorage != nil {
		if blobStorage.Tenant != nil {
			d.Set("tenant_id", blobStorage.Tenant)
		}

		if blobStorage.ServicePrincipalID != nil {
			d.Set("service_principal_id", blobStorage.ServicePrincipalID)
			d.Set("use_managed_identity", false)
		} else {
			d.Set("service_endpoint", blobStorage.ServiceEndpoint)
			d.Set("use_managed_identity", true)
		}
	}

	d.Set("additional_properties", blobStorage.AdditionalProperties)
	d.Set("description", blobStorage.Description)

	annotations := flattenDataFactoryAnnotations(blobStorage.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations` for Data Factory Linked Service Azure Blob Storage %q (Data Factory %q) / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	parameters := flattenDataFactoryParameters(blobStorage.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := blobStorage.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	return nil
}

func resourceDataFactoryLinkedServiceBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
			return fmt.Errorf("Error deleting Data Factory Linked Service BlobStorage %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
