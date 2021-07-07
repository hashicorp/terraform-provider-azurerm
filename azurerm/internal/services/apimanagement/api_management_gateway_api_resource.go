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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementGatewayApi() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayApiCreate,
		Read:   resourceApiManagementGatewayApiRead,
		Delete: resourceApiManagementGatewayApiDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GatewayApiID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceApiManagementGatewayApiCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayApisClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiID, err := parse.ApiID(d.Get("api_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_id`: %v", err)
	}

	gatewayID, err := parse.GatewayID(d.Get("gateway_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `gateway_id`: %v", err)
	}

	exists, err := client.GetEntityTag(ctx, gatewayID.ResourceGroup, gatewayID.ServiceName, gatewayID.Name, apiID.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(exists) {
			return fmt.Errorf("checking for present of existing API %q / Gateway %q (API Management Service %q / Resource Group %q): %+v", apiID.Name, gatewayID, gatewayID.ServiceName, gatewayID.ResourceGroup, err)
		}
	}

	if !utils.ResponseWasNotFound(exists) {
		// TODO: can we pull this from somewhere?
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId
		resourceId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/apis/%s", subscriptionId, gatewayID.ResourceGroup, gatewayID.ServiceName, gatewayID, apiID.Name)
		return tf.ImportAsExistsError("azurerm_api_management_gateway_api", resourceId)
	}

	params := &apimanagement.AssociationContract{}
	resp, err := client.CreateOrUpdate(ctx, gatewayID.ResourceGroup, gatewayID.ServiceName, gatewayID.Name, apiID.Name, params)
	if err != nil {
		return fmt.Errorf("adding API %q to Gateway %q (API Management Service %q / Resource Group %q): %+v", apiID.Name, gatewayID.Name, gatewayID.ServiceName, gatewayID.ResourceGroup, err)
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

	apiId := parse.NewApiID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName)

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

	gateway := parse.NewGatewayID(id.SubscriptionId, resourceGroup, serviceName, gatewayID)

	d.Set("api_id", apiId.ID())
	d.Set("gateway_id", gateway.ID())

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
