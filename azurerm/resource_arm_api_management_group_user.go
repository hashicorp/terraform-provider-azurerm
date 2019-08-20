package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementGroupUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementGroupUserCreate,
		Read:   resourceArmApiManagementGroupUserRead,
		Delete: resourceArmApiManagementGroupUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"user_id": azure.SchemaApiManagementChildName(),

			"group_name": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),
		},
	}
}

func resourceArmApiManagementGroupUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.GroupUsersClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	groupName := d.Get("group_name").(string)
	userId := d.Get("user_id").(string)

	if requireResourcesToBeImported {
		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, groupName, userId)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Error checking for present of existing User %q / Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp) {
			subscriptionId := meta.(*ArmClient).subscriptionId
			resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/groups/%s/users/%s", subscriptionId, resourceGroup, serviceName, groupName, userId)
			return tf.ImportAsExistsError("azurerm_api_management_group_user", resourceId)
		}
	}

	resp, err := client.Create(ctx, resourceGroup, serviceName, groupName, userId)
	if err != nil {
		return fmt.Errorf("Error adding User %q to Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
	}

	// there's no Read so this is best-effort
	d.SetId(*resp.ID)

	return resourceArmApiManagementGroupUserRead(d, meta)
}

func resourceArmApiManagementGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.GroupUsersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	groupName := id.Path["groups"]
	userId := id.Path["users"]

	resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, groupName, userId)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] User %q was not found in Group %q (API Management Service %q / Resource Group %q) was not found - removing from state!", userId, groupName, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving User %q / Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
	}

	d.Set("group_name", groupName)
	d.Set("user_id", userId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceArmApiManagementGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.GroupUsersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	groupName := id.Path["groups"]
	userId := id.Path["users"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, groupName, userId); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error removing User %q from Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
		}
	}

	return nil
}
