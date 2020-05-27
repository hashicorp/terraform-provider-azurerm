package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	groupName := d.Get("group_name").(string)
	userId := d.Get("user_id").(string)

	if features.ShouldResourcesBeImported() {
		resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, groupName, userId)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("checking for present of existing User %q / Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp) {
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/groups/%s/users/%s", subscriptionId, resourceGroup, serviceName, groupName, userId)
			return tf.ImportAsExistsError("azurerm_api_management_group_user", resourceId)
		}
	}

	resp, err := client.Create(ctx, resourceGroup, serviceName, groupName, userId)
	if err != nil {
		return fmt.Errorf("adding User %q to Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
	}

	// there's no Read so this is best-effort
	d.SetId(*resp.ID)

	return resourceArmApiManagementGroupUserRead(d, meta)
}

func resourceArmApiManagementGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

		return fmt.Errorf("retrieving User %q / Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
	}

	d.Set("group_name", groupName)
	d.Set("user_id", userId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceArmApiManagementGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
			return fmt.Errorf("removing User %q from Group %q (API Management Service %q / Resource Group %q): %+v", userId, groupName, serviceName, resourceGroup, err)
		}
	}

	return nil
}
