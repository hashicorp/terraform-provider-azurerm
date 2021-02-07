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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	keyVaultParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryLinkedServiceKeyVault() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataFactoryLinkedServiceKeyVaultCreateUpdate,
		Read:   resourceDataFactoryLinkedServiceKeyVaultRead,
		Update: resourceDataFactoryLinkedServiceKeyVaultCreateUpdate,
		Delete: resourceDataFactoryLinkedServiceKeyVaultDelete,

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

			"key_vault_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
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

func resourceDataFactoryLinkedServiceKeyVaultCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	keyVaultId, err := keyVaultParse.VaultID(d.Get("key_vault_id").(string))
	if err != nil {
		return err
	}

	keyVaultBaseUri, err := keyVaultsClient.BaseUriForKeyVault(ctx, *keyVaultId)
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_linked_service_key_vault", *existing.ID)
		}
	}

	azureKeyVaultProperties := &datafactory.AzureKeyVaultLinkedServiceTypeProperties{
		BaseURL: keyVaultBaseUri,
	}

	azureKeyVaultLinkedService := &datafactory.AzureKeyVaultLinkedService{
		Description:                              utils.String(d.Get("description").(string)),
		AzureKeyVaultLinkedServiceTypeProperties: azureKeyVaultProperties,
		Type:                                     datafactory.TypeAzureKeyVault,
	}

	if v, ok := d.GetOk("parameters"); ok {
		azureKeyVaultLinkedService.Parameters = expandDataFactoryParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("integration_runtime_name"); ok {
		azureKeyVaultLinkedService.ConnectVia = expandDataFactoryLinkedServiceIntegrationRuntime(v.(string))
	}

	if v, ok := d.GetOk("additional_properties"); ok {
		azureKeyVaultLinkedService.AdditionalProperties = v.(map[string]interface{})
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		azureKeyVaultLinkedService.Annotations = &annotations
	}

	linkedService := datafactory.LinkedServiceResource{
		Properties: azureKeyVaultLinkedService,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, dataFactoryName, name, linkedService, ""); err != nil {
		return fmt.Errorf("Error creating/updating Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceDataFactoryLinkedServiceKeyVaultRead(d, meta)
}

func resourceDataFactoryLinkedServiceKeyVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
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

		return fmt.Errorf("Error retrieving Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	keyVault, ok := resp.Properties.AsAzureKeyVaultLinkedService()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", name, dataFactoryName, resourceGroup, datafactory.TypeAzureKeyVault, *resp.Type)
	}

	d.Set("additional_properties", keyVault.AdditionalProperties)
	d.Set("description", keyVault.Description)

	annotations := flattenDataFactoryAnnotations(keyVault.Annotations)
	if err := d.Set("annotations", annotations); err != nil {
		return fmt.Errorf("Error setting `annotations`: %+v", err)
	}

	parameters := flattenDataFactoryParameters(keyVault.Parameters)
	if err := d.Set("parameters", parameters); err != nil {
		return fmt.Errorf("Error setting `parameters`: %+v", err)
	}

	if connectVia := keyVault.ConnectVia; connectVia != nil {
		if connectVia.ReferenceName != nil {
			d.Set("integration_runtime_name", connectVia.ReferenceName)
		}
	}

	baseUrl := ""
	if properties := keyVault.AzureKeyVaultLinkedServiceTypeProperties; properties != nil {
		if properties.BaseURL != nil {
			val, ok := properties.BaseURL.(string)
			if ok {
				baseUrl = val
			} else {
				log.Printf("[DEBUG] Skipping base url string %q since it's not a string", val)
			}
		}
	}

	var keyVaultId *string
	if baseUrl != "" {
		keyVaultId, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, baseUrl)
		if err != nil {
			return err
		}
	}
	d.Set("key_vault_id", keyVaultId)

	return nil
}

func resourceDataFactoryLinkedServiceKeyVaultDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Data Factory Linked Service Key Vault %q (Data Factory %q / Resource Group %q): %+v", name, dataFactoryName, resourceGroup, err)
		}
	}

	return nil
}
