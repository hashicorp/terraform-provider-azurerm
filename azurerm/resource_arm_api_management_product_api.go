package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementProductApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementProductApiCreate,
		Read:   resourceArmApiManagementProductApiRead,
		Delete: resourceArmApiManagementProductApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"api_name": azure.SchemaApiManagementChildName(),

			"product_id": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),
		},
	}
}

func resourceArmApiManagementProductApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ProductApisClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)
	productId := d.Get("product_id").(string)

	if features.ShouldResourcesBeImported() {
		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, apiName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Error checking for present of existing API %q / Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp) {
			subscriptionId := meta.(*ArmClient).subscriptionId
			resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/apis/%s", subscriptionId, resourceGroup, serviceName, productId, apiName)
			return tf.ImportAsExistsError("azurerm_api_management_product_api", resourceId)
		}
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, productId, apiName)
	if err != nil {
		return fmt.Errorf("Error adding API %q to Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
	}

	// there's no Read so this is best-effort
	d.SetId(*resp.ID)

	return resourceArmApiManagementProductApiRead(d, meta)
}

func resourceArmApiManagementProductApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ProductApisClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	productId := id.Path["products"]
	apiName := id.Path["apis"]

	resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, apiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] API %q was not found in Product  %q (API Management Service %q / Resource Group %q) was not found - removing from state!", apiName, productId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving API %q / Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
	}

	d.Set("api_name", apiName)
	d.Set("product_id", productId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceArmApiManagementProductApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ProductApisClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	productId := id.Path["products"]
	apiName := id.Path["apis"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, productId, apiName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error removing API %q from Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
