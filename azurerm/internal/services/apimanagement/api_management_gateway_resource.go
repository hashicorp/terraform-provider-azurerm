package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayCreateUpdate,
		Read:   resourceApiManagementGatewayRead,
		Update: resourceApiManagementGatewayCreateUpdate,
		Delete: resourceApiManagementGatewayDelete,

		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(10 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"gateway_id": schemaz.SchemaApiManagementChildName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 2000),
			},
		},
	}
}

func resourceApiManagementGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("gateway_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing gateway %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_gateway", *existing.ID)
		}
	}
	location := d.Get("location").(string)
	locationData := &apimanagement.ResourceLocationDataContract{
		Name: &location,
	}
	gatewayContract := apimanagement.GatewayContract{
		GatewayContractProperties: &apimanagement.GatewayContractProperties{
			LocationData: locationData,
		},
	}
	if description, ok := d.GetOk("description"); ok {
		gatewayContract.GatewayContractProperties.Description = utils.String(description.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, gatewayContract, ""); err != nil {
		return fmt.Errorf("creating/updating gateway %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("retrieving gateway %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for gateway %q (API Management Service %q / Resource Group %q)", name, serviceName, resourceGroup)
	}

	d.SetId(*read.ID)
	return resourceApiManagementGatewayRead(d, meta)
}

func resourceApiManagementGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.GatewayID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	name := id.Name

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] gateway %q (API Management Service %q / Resource Group %q) does not exist - removing from state!", name, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving gateway %q (API Management Service %q / Resource Group %q): %+v", name, serviceName, resourceGroup, err)
	}

	d.Set("gateway_id", name)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.GatewayContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("location", string(*props.LocationData.Name))
	}

	return nil
}

func resourceApiManagementGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	name := id.Name

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting gateway %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
		}
	}

	return nil
}
