package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementProductApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementProductApiCreate,
		Read:   resourceApiManagementProductApiRead,
		Delete: resourceApiManagementProductApiDelete,
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
			"api_name": azure.SchemaApiManagementApiName(),

			"product_id": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementProductApiCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)
	productId := d.Get("product_id").(string)

	exists, err := client.CheckEntityExists(ctx, resourceGroup, serviceName, productId, apiName)
	if err != nil {
		if !utils.ResponseWasNotFound(exists) {
			return fmt.Errorf("checking for present of existing API %q / Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(exists) {
		// TODO: can we pull this from somewhere?
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/products/%s/apis/%s", subscriptionId, resourceGroup, serviceName, productId, apiName)
		return tf.ImportAsExistsError("azurerm_api_management_product_api", resourceId)
	}

	resp, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, productId, apiName)
	if err != nil {
		return fmt.Errorf("adding API %q to Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
	}

	// there's no Read so this is best-effort
	d.SetId(*resp.ID)

	return resourceApiManagementProductApiRead(d, meta)
}

func resourceApiManagementProductApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

		return fmt.Errorf("retrieving API %q / Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
	}

	// This can be removed once updated to apimanagement API to 2019-01-01
	// https://github.com/Azure/azure-sdk-for-go/blob/master/services/apimanagement/mgmt/2019-01-01/apimanagement/productapi.go#L134
	if utils.ResponseWasNotFound(resp) {
		log.Printf("[DEBUG] API %q was not found in Product  %q (API Management Service %q / Resource Group %q) was not found - removing from state!", apiName, productId, serviceName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("api_name", apiName)
	d.Set("product_id", productId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceApiManagementProductApiDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductApisClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
			return fmt.Errorf("removing API %q from Product %q (API Management Service %q / Resource Group %q): %+v", apiName, productId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
