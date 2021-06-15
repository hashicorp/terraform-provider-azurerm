package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceExpressRouteCircuitConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceExpressRouteCircuitConnectionCreateUpdate,
		Read:   resourceExpressRouteCircuitConnectionRead,
		Update: resourceExpressRouteCircuitConnectionCreateUpdate,
		Delete: resourceExpressRouteCircuitConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ExpressRouteCircuitConnectionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitConnectionName,
			},

			"peering_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitPeeringID,
			},

			"peer_peering_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitPeeringID,
			},

			"address_prefix_ipv4": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},

			"authorization_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.IsUUID,
			},

			"address_prefix_ipv6": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
			},
		},
	}
}

func resourceExpressRouteCircuitConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	expressRouteCircuitPeeringId, err := parse.ExpressRouteCircuitPeeringID(d.Get("peering_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewExpressRouteCircuitConnectionID(expressRouteCircuitPeeringId.SubscriptionId, expressRouteCircuitPeeringId.ResourceGroup, expressRouteCircuitPeeringId.ExpressRouteCircuitName, string(network.ExpressRoutePeeringTypeAzurePrivatePeering), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, string(network.ExpressRoutePeeringTypeAzurePrivatePeering), id.ConnectionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_express_route_circuit_connection", id.ID())
		}
	}

	expressRouteCircuitConnectionParameters := network.ExpressRouteCircuitConnection{
		Name: utils.String(id.ConnectionName),
		ExpressRouteCircuitConnectionPropertiesFormat: &network.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: utils.String(d.Get("address_prefix_ipv4").(string)),
			ExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(expressRouteCircuitPeeringId.ID()),
			},
			PeerExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(d.Get("peer_peering_id").(string)),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.AuthorizationKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("address_prefix_ipv6"); ok {
		expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.Ipv6CircuitConnectionConfig = &network.Ipv6CircuitConnectionConfig{
			AddressPrefix: utils.String(v.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, string(network.ExpressRoutePeeringTypeAzurePrivatePeering), id.ConnectionName, expressRouteCircuitConnectionParameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitConnectionRead(d, meta)
}

func resourceExpressRouteCircuitConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("peering_id", parse.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, string(network.ExpressRoutePeeringTypeAzurePrivatePeering)).ID())

	if props := resp.ExpressRouteCircuitConnectionPropertiesFormat; props != nil {
		d.Set("address_prefix_ipv4", props.AddressPrefix)

		if props.AuthorizationKey != nil {
			d.Set("authorization_key", props.AuthorizationKey)
		} else {
			d.Set("authorization_key", "")
		}

		if props.Ipv6CircuitConnectionConfig != nil && props.Ipv6CircuitConnectionConfig.AddressPrefix != nil {
			d.Set("address_prefix_ipv6", props.Ipv6CircuitConnectionConfig.AddressPrefix)
		} else {
			d.Set("address_prefix_ipv6", "")
		}

		if props.PeerExpressRouteCircuitPeering != nil && props.PeerExpressRouteCircuitPeering.ID != nil {
			d.Set("peer_peering_id", resp.PeerExpressRouteCircuitPeering.ID)
		}
	}

	return nil
}

func resourceExpressRouteCircuitConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
