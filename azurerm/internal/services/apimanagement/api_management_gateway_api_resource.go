package apimanagement

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementGatewayApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayApiCreate,
		Read:   resourceApiManagementGatewayApiRead,
		Delete: resourceApiManagementGatewayApiDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(15 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(15 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_name": schemaz.SchemaApiManagementApiName(),

			"gateway_id": schemaz.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),
		},
	}
}

func resourceApiManagementGatewayApiCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)
	gatewayID := d.Get("gateway_id").(string)

	exists, err := client.GetEntityTag(ctx, resourceGroup, serviceName, gatewayID, apiName)
	if err != nil {
		if !utils.ResponseWasNotFound(exists) {
			return fmt.Errorf("checking for present of existing API %q / Gateway %q (API Management Service %q / Resource Group %q): %+v", apiName, gatewayID, serviceName, resourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(exists) {
		// TODO: can we pull this from somewhere?
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/apis/%s", subscriptionId, resourceGroup, serviceName, gatewayID, apiName)
		return tf.ImportAsExistsError("azurerm_api_management_gateway_api", resourceId)
	}

	params := &apimanagement.AssociationContract{}
	resp, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, gatewayID, apiName, params)
	if err != nil {
		return fmt.Errorf("adding API %q to Gateway %q (API Management Service %q / Resource Group %q): %+v", apiName, gatewayID, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceApiManagementGatewayApiRead(d, meta)
}

func resourceApiManagementGatewayApiRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayApiID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	gatewayID := id.GatewayName
	apiName := id.ApiName

	resp, err := client.GetEntityTag(ctx, resourceGroup, serviceName, gatewayID, apiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] gateway %q (API Management Service %q / Resource Group %q) does not exist - removing from state!", apiName, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		if utils.ResponseWasStatusCode(resp, http.StatusNoContent) {
			log.Printf("[DEBUG] gateway %q (API Management Service %q / Resource Group %q) returning a no content status - bypassing and moving on!", apiName, serviceName, resourceGroup)
		} else {
			return fmt.Errorf("retrieving gateway %q (API Management Service %q / Resource Group %q): %+v", apiName, serviceName, resourceGroup, err)
		}
	}

	if utils.ResponseWasNotFound(resp) {
		log.Printf("[DEBUG] API %q was not found in Gateway  %q (API Management Service %q / Resource Group %q) was not found - removing from state!", apiName, gatewayID, serviceName, resourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("api_name", apiName)
	d.Set("gateway_id", gatewayID)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	return nil
}

func resourceApiManagementGatewayApiDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayApiID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	gatewayID := id.GatewayName
	apiName := id.ApiName

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, gatewayID, apiName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing API %q from Gateway %q (API Management Service %q / Resource Group %q): %+v", apiName, gatewayID, serviceName, resourceGroup, err)
		}
	}

	return nil
}
