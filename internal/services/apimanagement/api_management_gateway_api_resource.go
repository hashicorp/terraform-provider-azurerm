package apimanagement

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
		if !utils.ResponseWasStatusCode(exists, http.StatusNoContent) {
			if !utils.ResponseWasNotFound(exists) {
				return fmt.Errorf("checking for presence of existing %s: %+v", gatewayID, err)
			}
		}
	}

	id := parse.NewGatewayApiID(gatewayID.SubscriptionId, gatewayID.ResourceGroup, gatewayID.ServiceName, gatewayID.Name, apiID.Name)
	if !utils.ResponseWasNotFound(exists) {
		return tf.ImportAsExistsError("azurerm_api_management_gateway_api", id.ID())
	}
	params := &apimanagement.AssociationContract{}
	_, err = client.CreateOrUpdate(ctx, gatewayID.ResourceGroup, gatewayID.ServiceName, gatewayID.Name, apiID.Name, params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

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

	apiId := parse.NewApiID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName)
	resp, err := client.GetEntityTag(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.ApiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", id)
			d.SetId("")
			return nil
		}
		if utils.ResponseWasStatusCode(resp, http.StatusNoContent) {
			log.Printf("[DEBUG] %s returned with No Content status - bypassing and moving on!", id)
		} else {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
	}
	if utils.ResponseWasNotFound(resp) {
		log.Printf("[DEBUG] %s was not found - removing from state!", id)
		d.SetId("")
		return nil
	}
	gateway := parse.NewGatewayID(id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.GatewayName)

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

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.ApiName); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing %s: %+v", id, err)
		}
	}

	return nil
}
