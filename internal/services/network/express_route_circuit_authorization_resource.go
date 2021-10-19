package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceExpressRouteCircuitAuthorization() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitAuthorizationCreate,
		Read:   resourceExpressRouteCircuitAuthorizationRead,
		Delete: resourceExpressRouteCircuitAuthorizationDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AuthorizationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"express_route_circuit_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"authorization_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"authorization_use_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceExpressRouteCircuitAuthorizationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteAuthsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	circuitName := d.Get("express_route_circuit_name").(string)

	locks.ByName(circuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(circuitName, expressRouteCircuitResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, circuitName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %s", name, circuitName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_authorization", *existing.ID)
		}
	}

	properties := network.ExpressRouteCircuitAuthorization{
		AuthorizationPropertiesFormat: &network.AuthorizationPropertiesFormat{},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, circuitName, name, properties)
	if err != nil {
		return fmt.Errorf("Creating/Updating Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Express Route Circuit Authorization %q (Circuit %q / Resource Group %q) to finish creating/updating: %+v", name, circuitName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, circuitName, name)
	if err != nil {
		return fmt.Errorf("retrieving Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", name, circuitName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceExpressRouteCircuitAuthorizationRead(d, meta)
}

func resourceExpressRouteCircuitAuthorizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteAuthsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", id.Name, id.ExpressRouteCircuitName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("express_route_circuit_name", id.ExpressRouteCircuitName)

	if props := resp.AuthorizationPropertiesFormat; props != nil {
		d.Set("authorization_key", props.AuthorizationKey)
		d.Set("authorization_use_status", string(props.AuthorizationUseStatus))
	}

	return nil
}

func resourceExpressRouteCircuitAuthorizationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteAuthsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AuthorizationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Express Route Circuit Authorization %q (Circuit %q / Resource Group %q): %+v", id.Name, id.ExpressRouteCircuitName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for Express Route Circuit Authorization %q (Circuit %q / Resource Group %q) to be deleted: %+v", id.Name, id.ExpressRouteCircuitName, id.ResourceGroup, err)
	}

	return nil
}
