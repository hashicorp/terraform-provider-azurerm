package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementApiSchema() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementApiSchemaCreateUpdate,
		Read:   resourceArmApiManagementApiSchemaRead,
		Update: resourceArmApiManagementApiSchemaCreateUpdate,
		Delete: resourceArmApiManagementApiSchemaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"schema_id": azure.SchemaApiManagementChildName(),

			"api_name": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"content_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmApiManagementApiSchemaCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ApiSchemasClient
	ctx := meta.(*ArmClient).StopContext

	schemaID := d.Get("schema_id").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_schema", *existing.ID)
		}
	}

	contentType := d.Get("content_type").(string)
	value := d.Get("value").(string)
	parameters := apimanagement.SchemaContract{
		SchemaContractProperties: &apimanagement.SchemaContractProperties{
			ContentType: &contentType,
			SchemaDocumentProperties: &apimanagement.SchemaDocumentProperties{
				Value: &value,
			},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, schemaID, parameters, ""); err != nil {
		return fmt.Errorf("Error creating or updating API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
	if err != nil {
		return fmt.Errorf("Error retrieving API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementApiSchemaRead(d, meta)
}

func resourceArmApiManagementApiSchemaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ApiSchemasClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiName := id.Path["apis"]
	schemaID := id.Path["schemas"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Schema %q (API Management Service %q / API %q / Resource Group %q) was not found - removing from state!", schemaID, serviceName, apiName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("api_name", apiName)
	d.Set("schema_id", schemaID)

	if properties := resp.SchemaContractProperties; properties != nil {
		d.Set("content_type", properties.ContentType)
		if documentProperties := properties.SchemaDocumentProperties; documentProperties != nil {
			d.Set("value", documentProperties.Value)
		}
	}

	return nil
}

func resourceArmApiManagementApiSchemaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ApiSchemasClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiName := id.Path["apis"]
	schemaID := id.Path["schemas"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, apiName, schemaID, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
		}
	}

	return nil
}
