package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceExpressRouteGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceExpressRouteGatewayCreateUpdate,
		Read:   resourceExpressRouteGatewayRead,
		Update: resourceExpressRouteGatewayCreateUpdate,
		Delete: resourceExpressRouteGatewayDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ExpressRouteGatewayID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"scale_units": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 10),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceExpressRouteGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.ExpressRouteGatewaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Println("[INFO] preparing arguments for ExpressRoute Gateway creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewExpressRouteGatewayID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing ExpressRoute Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_express_route_gateway", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	virtualHubId := d.Get("virtual_hub_id").(string)
	t := d.Get("tags").(map[string]interface{})

	minScaleUnits := int32(d.Get("scale_units").(int))
	parameters := network.ExpressRouteGateway{
		Location: utils.String(location),
		ExpressRouteGatewayProperties: &network.ExpressRouteGatewayProperties{
			AutoScaleConfiguration: &network.ExpressRouteGatewayPropertiesAutoScaleConfiguration{
				Bounds: &network.ExpressRouteGatewayPropertiesAutoScaleConfigurationBounds{
					Min: &minScaleUnits,
				},
			},
			VirtualHub: &network.VirtualHubID{
				ID: &virtualHubId,
			},
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating ExpressRoute Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of ExpressRoute Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteGatewayRead(d, meta)
}

func resourceExpressRouteGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] ExpressRoute Gateway %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading ExpressRoute Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ExpressRouteGatewayProperties; props != nil {
		virtualHubId := ""
		if props.VirtualHub != nil && props.VirtualHub.ID != nil {
			virtualHubId = *props.VirtualHub.ID
		}
		d.Set("virtual_hub_id", virtualHubId)

		scaleUnits := 0
		if props.AutoScaleConfiguration != nil && props.AutoScaleConfiguration.Bounds != nil && props.AutoScaleConfiguration.Bounds.Min != nil {
			scaleUnits = int(*props.AutoScaleConfiguration.Bounds.Min)
		}
		d.Set("scale_units", scaleUnits)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceExpressRouteGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteGatewayID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting ExpressRoute Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting ExpressRoute Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
