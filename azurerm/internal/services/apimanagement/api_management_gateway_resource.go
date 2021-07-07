package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayCreateUpdate,
		Read:   resourceApiManagementGatewayRead,
		Update: resourceApiManagementGatewayCreateUpdate,
		Delete: resourceApiManagementGatewayDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApiManagementID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"location_data": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"city": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"district": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"region": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceApiManagementGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := parse.ApiManagementID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := parse.NewGatewayID(apimId.SubscriptionId, apimId.ResourceGroup, apimId.ServiceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("making read request %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_gateway", id.ID())
		}
	}

	description := d.Get("description").(string)
	locationData := expandApiManagementGatewayLocationData(d.Get("location_data").([]interface{}))

	parameters := apimanagement.GatewayContract{
		GatewayContractProperties: &apimanagement.GatewayContractProperties{
			Description:  utils.String(description),
			LocationData: locationData,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.Name, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
	apimId := parse.NewApiManagementID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Gateway %q (Resource Group %q / API Management Service %q) was not found - removing from state!", name, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("api_management_id", apimId.ID())

	if properties := resp.GatewayContractProperties; properties != nil {
		d.Set("description", properties.Description)
		d.Set("location_data", flattenApiManagementGatewayLocationData(properties.LocationData))
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
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func expandApiManagementGatewayLocationData(input []interface{}) *apimanagement.ResourceLocationDataContract {
	if len(input) == 0 {
		return nil
	}

	locationData := apimanagement.ResourceLocationDataContract{}

	vs := input[0].(map[string]interface{})
	for k, v := range vs {
		switch k {
		case "name":
			locationData.Name = utils.String(v.(string))
		case "city":
			locationData.City = utils.String(v.(string))
		case "district":
			locationData.District = utils.String(v.(string))
		case "region":
			locationData.CountryOrRegion = utils.String(v.(string))
		}
	}

	return &locationData
}

func flattenApiManagementGatewayLocationData(input *apimanagement.ResourceLocationDataContract) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	locationData := map[string]interface{}{
		"name":     utils.NormalizeNilableString(input.Name),
		"city":     utils.NormalizeNilableString(input.City),
		"region":   utils.NormalizeNilableString(input.CountryOrRegion),
		"district": utils.NormalizeNilableString(input.District),
	}

	return []interface{}{locationData}
}
