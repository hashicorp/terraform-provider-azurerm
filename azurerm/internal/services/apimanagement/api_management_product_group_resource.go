package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementProductGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductGroupCreate,
		Read:   resourceApiManagementProductGroupRead,
		Delete: resourceApiManagementProductGroupDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"product_id": schemaz.SchemaApiManagementChildName(),

			"group_name": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementProductGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	groupName := d.Get("group_name").(string)
	productId := d.Get("product_id").(string)

	exists, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, groupName)
	if err != nil {
		if !utils.ResponseWasNotFound(exists) {
			return fmt.Errorf("checking for present of existing Product %q / Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(exists) {
		// TODO: can we pull this from somewhere?
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/groups/%s", subscriptionId, resourceGroup, serviceName, productId, groupName)
		return tf.ImportAsExistsError("azurerm_api_management_product_group", resourceId)
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, productId, groupName)
	if err != nil {
		return fmt.Errorf("adding Product %q to Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceApiManagementProductGroupRead(d, meta)
}

func resourceApiManagementProductGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductGroupID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	groupName := id.GroupName
	productId := id.ProductName

	resp, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, groupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] Product %q was not found in Group %q (API Management Service %q / Resource Group %q) was not found - removing from state!", productId, groupName, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Product %q / Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
	}

	d.Set("group_name", groupName)
	d.Set("product_id", productId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceApiManagementProductGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductGroupID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	groupName := id.GroupName
	productId := id.ProductName

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, productId, groupName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing Product %q from Group %q (API Management Service %q / Resource Group %q): %+v", productId, groupName, serviceName, resourceGroup, err)
		}
	}

	return nil
}
