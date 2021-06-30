package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceExpressRouteCircuitConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitConnectionCreate,
		Read:   resourceExpressRouteCircuitConnectionRead,
		Update: resourceExpressRouteCircuitConnectionUpdate,
		Delete: resourceExpressRouteCircuitConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRouteCircuitConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitConnectionName,
			},

			"peering_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitPeeringID,
			},

			"peer_peering_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitPeeringID,
			},

			"address_prefix_ipv4": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},

			"authorization_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.IsUUID,
			},

			"address_prefix_ipv6": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
			},
		},
	}
}

func resourceExpressRouteCircuitConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	circuitClient := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	circuitPeeringId, err := parse.ExpressRouteCircuitPeeringID(d.Get("peering_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewExpressRouteCircuitConnectionID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroup, circuitPeeringId.ExpressRouteCircuitName, circuitPeeringId.PeeringName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_express_route_circuit_connection", id.ID())
	}

	circuitPeerPeeringId, err := parse.ExpressRouteCircuitPeeringID(d.Get("peer_peering_id").(string))
	if err != nil {
		return err
	}

	expressRouteCircuitConnectionParameters := network.ExpressRouteCircuitConnection{
		Name: utils.String(id.ConnectionName),
		ExpressRouteCircuitConnectionPropertiesFormat: &network.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: utils.String(d.Get("address_prefix_ipv4").(string)),
			ExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(circuitPeeringId.ID()),
			},
			PeerExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(circuitPeerPeeringId.ID()),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.AuthorizationKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("address_prefix_ipv6"); ok {
		circuitId := parse.NewExpressRouteCircuitID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroup, circuitPeeringId.ExpressRouteCircuitName)

		circuit, err := circuitClient.Get(ctx, circuitId.ResourceGroup, circuitId.Name)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", circuitId, err)
		}

		if circuit.ExpressRouteCircuitPropertiesFormat != nil && circuit.ExpressRouteCircuitPropertiesFormat.ExpressRoutePort != nil {
			return fmt.Errorf("`address_prefix_ipv6` cannot be set when ExpressRoute Circuit Connection with ExpressRoute Circuit based on ExpressRoute Port")
		} else {
			expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.Ipv6CircuitConnectionConfig = &network.Ipv6CircuitConnectionConfig{
				AddressPrefix: utils.String(v.(string)),
			}
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName, expressRouteCircuitConnectionParameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitConnectionRead(d, meta)
}

func resourceExpressRouteCircuitConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("peering_id", parse.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName).ID())

	if props := resp.ExpressRouteCircuitConnectionPropertiesFormat; props != nil {
		d.Set("address_prefix_ipv4", props.AddressPrefix)

		// The ExpressRoute Circuit Connection API returns "*****************" for AuthorizationKey when it's changed from a valid value to `nil`
		// See more details from https://github.com/Azure/azure-rest-api-specs/issues/15030
		authorizationKey := ""
		if props.AuthorizationKey != nil && *props.AuthorizationKey != "*****************" {
			authorizationKey = *props.AuthorizationKey
		}
		d.Set("authorization_key", authorizationKey)

		addressPrefixIPv6 := ""
		if props.Ipv6CircuitConnectionConfig != nil && props.Ipv6CircuitConnectionConfig.AddressPrefix != nil {
			addressPrefixIPv6 = *props.Ipv6CircuitConnectionConfig.AddressPrefix
		}
		d.Set("address_prefix_ipv6", addressPrefixIPv6)

		if props.PeerExpressRouteCircuitPeering != nil && props.PeerExpressRouteCircuitPeering.ID != nil {
			circuitPeerPeeringId, err := parse.ExpressRouteCircuitPeeringID(*props.PeerExpressRouteCircuitPeering.ID)
			if err != nil {
				return err
			}
			d.Set("peer_peering_id", circuitPeerPeeringId.ID())
		}
	}

	return nil
}

func resourceExpressRouteCircuitConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	circuitClient := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitConnectionID(d.Id())
	if err != nil {
		return err
	}

	circuitPeeringId, err := parse.ExpressRouteCircuitPeeringID(d.Get("peering_id").(string))
	if err != nil {
		return err
	}

	circuitPeerPeeringId, err := parse.ExpressRouteCircuitPeeringID(d.Get("peer_peering_id").(string))
	if err != nil {
		return err
	}

	expressRouteCircuitConnectionParameters := network.ExpressRouteCircuitConnection{
		Name: utils.String(id.ConnectionName),
		ExpressRouteCircuitConnectionPropertiesFormat: &network.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: utils.String(d.Get("address_prefix_ipv4").(string)),
			ExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(circuitPeeringId.ID()),
			},
			PeerExpressRouteCircuitPeering: &network.SubResource{
				ID: utils.String(circuitPeerPeeringId.ID()),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.AuthorizationKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("address_prefix_ipv6"); ok {
		circuitId := parse.NewExpressRouteCircuitID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroup, circuitPeeringId.ExpressRouteCircuitName)

		circuit, err := circuitClient.Get(ctx, circuitId.ResourceGroup, circuitId.Name)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", circuitId, err)
		}

		if circuit.ExpressRouteCircuitPropertiesFormat != nil && circuit.ExpressRouteCircuitPropertiesFormat.ExpressRoutePort != nil {
			return fmt.Errorf("`address_prefix_ipv6` cannot be set when ExpressRoute Circuit Connection with ExpressRoute Circuit based on ExpressRoute Port")
		} else {
			expressRouteCircuitConnectionParameters.ExpressRouteCircuitConnectionPropertiesFormat.Ipv6CircuitConnectionConfig = &network.Ipv6CircuitConnectionConfig{
				AddressPrefix: utils.String(v.(string)),
			}
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ExpressRouteCircuitName, id.PeeringName, id.ConnectionName, expressRouteCircuitConnectionParameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceExpressRouteCircuitConnectionRead(d, meta)
}

func resourceExpressRouteCircuitConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
