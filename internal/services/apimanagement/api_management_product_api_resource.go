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

func resourceApiManagementProductApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductApiCreate,
		Read:   resourceApiManagementProductApiRead,
		Delete: resourceApiManagementProductApiDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ProductApiID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_name": schemaz.SchemaApiManagementApiName(),

			"product_id": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementProductApiCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewProductApiID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string), d.Get("api_name").(string))

	exists, err := client.CheckEntityExists(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.ApiName)
	if err != nil {
		if !utils.ResponseWasNotFound(exists) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(exists) {
		return tf.ImportAsExistsError("azurerm_api_management_product_api", id.ID())
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.ApiName); err != nil {
		return fmt.Errorf("adding API %q to Product %q (API Management Service %q / Resource Group %q): %+v", id.ApiName, id.ProductName, id.ServiceName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementProductApiRead(d, meta)
}

func resourceApiManagementProductApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductApiID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CheckEntityExists(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.ApiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] API %q was not found in Product  %q (API Management Service %q / Resource Group %q) was not found - removing from state!", id.ApiName, id.ProductName, id.ServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_name", id.ApiName)
	d.Set("product_id", id.ProductName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	return nil
}

func resourceApiManagementProductApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductApiID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.ProductName, id.ApiName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing API %q from Product %q (API Management Service %q / Resource Group %q): %+v", id.ApiName, id.ProductName, id.ServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
