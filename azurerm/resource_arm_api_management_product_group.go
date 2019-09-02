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

func resourceArmApiManagementProductGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementProductGroupCreate,
		Read:   resourceArmApiManagementProductGroupRead,
		Delete: resourceArmApiManagementProductGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"product_id": azure.SchemaApiManagementChildName(),

			"group_name": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),
		},
	}
}

func resourceArmApiManagementProductGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ProductGroupsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	groupName := d.Get("group_name").(string)
	productId := d.Get("product_id").(string)

	if features.ShouldResourcesBeImported() {
		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, groupName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Error checking for present of existing Product %q / Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp) {
			subscriptionId := meta.(*ArmClient).subscriptionId
			resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/groups/%s", subscriptionId, resourceGroup, serviceName, groupName, productId)
			return tf.ImportAsExistsError("azurerm_api_management_product_group", resourceId)
		}
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, productId, groupName)
	if err != nil {
		return fmt.Errorf("Error adding Product %q to Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
	}

	// there's no Read so this is best-effort
	d.SetId(*resp.ID)

	return resourceArmApiManagementProductGroupRead(d, meta)
}

func resourceArmApiManagementProductGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ProductGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	groupName := id.Path["groups"]
	productId := id.Path["products"]

	resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, groupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Product %q was not found in Group %q (API Management Service %q / Resource Group %q) was not found - removing from state!", productId, groupName, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Product %q / Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
	}

	d.Set("group_name", groupName)
	d.Set("product_id", productId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceArmApiManagementProductGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.ProductGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	groupName := id.Path["groups"]
	productId := id.Path["products"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, productId, groupName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error removing Product %q from Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
		}
	}

	return nil
}
