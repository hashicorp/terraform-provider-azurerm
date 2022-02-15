package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementGroupUser() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGroupUserCreate,
		Read:   resourceApiManagementGroupUserRead,
		Delete: resourceApiManagementGroupUserDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GroupUserID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"user_id": schemaz.SchemaApiManagementChildName(),

			"group_name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementGroupUserCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewGroupUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("group_name").(string), d.Get("user_id").(string))

	exists, err := client.CheckEntityExists(ctx, id.ResourceGroup, id.ServiceName, id.GroupName, id.UserName)
	if err != nil {
		if !utils.ResponseWasNotFound(exists) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(exists) {
		return tf.ImportAsExistsError("azurerm_api_management_group_user", id.ID())
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.ServiceName, id.GroupName, id.UserName); err != nil {
		return fmt.Errorf("adding User %q to Group %q (API Management Service %q / Resource Group %q): %+v", id.UserName, id.GroupName, id.ServiceName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGroupUserRead(d, meta)
}

func resourceApiManagementGroupUserRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GroupUserID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CheckEntityExists(ctx, id.ResourceGroup, id.ServiceName, id.GroupName, id.UserName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("group_name", id.GroupName)
	d.Set("user_id", id.UserName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	return nil
}

func resourceApiManagementGroupUserDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupUsersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GroupUserID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.GroupName, id.UserName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing User %q from Group %q (API Management Service %q / Resource Group %q): %+v", id.UserName, id.GroupName, id.ServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
