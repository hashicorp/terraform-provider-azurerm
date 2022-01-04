package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
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
			_, err := parse.ExpressRouteCircuitAuthorizationID(id)
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := parse.NewExpressRouteCircuitAuthorizationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("express_route_circuit_name").(string), d.Get("name").(string))

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.AuthorizationName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_authorization", *existing.ID)
		}
	}

	properties := network.ExpressRouteCircuitAuthorization{
		AuthorizationPropertiesFormat: &network.AuthorizationPropertiesFormat{},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.AuthorizationName, properties)
	if err != nil {
		return fmt.Errorf("Creating/Updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for  %s to finish creating/updating: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitAuthorizationRead(d, meta)
}

func resourceExpressRouteCircuitAuthorizationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteAuthsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.AuthorizationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AuthorizationName)
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

	id, err := parse.ExpressRouteCircuitAuthorizationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.AuthorizationName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}
