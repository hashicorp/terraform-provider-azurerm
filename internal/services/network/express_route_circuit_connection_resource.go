// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuitconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			_, err := expressroutecircuitconnections.ParsePeeringConnectionID(id)
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
				ValidateFunc: commonids.ValidateExpressRouteCircuitPeeringID,
			},

			"peer_peering_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateExpressRouteCircuitPeeringID,
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
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnections
	circuitClient := meta.(*clients.Client).Network.ExpressRouteCircuits
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	circuitPeeringId, err := commonids.ParseExpressRouteCircuitPeeringID(d.Get("peering_id").(string))
	if err != nil {
		return err
	}

	id := expressroutecircuitconnections.NewPeeringConnectionID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroupName, circuitPeeringId.CircuitName, circuitPeeringId.PeeringName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_express_route_circuit_connection", id.ID())
	}

	circuitPeerPeeringId, err := commonids.ParseExpressRouteCircuitPeeringID(d.Get("peer_peering_id").(string))
	if err != nil {
		return err
	}

	expressRouteCircuitConnectionParameters := expressroutecircuitconnections.ExpressRouteCircuitConnection{
		Name: pointer.To(id.ConnectionName),
		Properties: &expressroutecircuitconnections.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: pointer.To(d.Get("address_prefix_ipv4").(string)),
			ExpressRouteCircuitPeering: &expressroutecircuitconnections.SubResource{
				Id: pointer.To(circuitPeeringId.ID()),
			},
			PeerExpressRouteCircuitPeering: &expressroutecircuitconnections.SubResource{
				Id: pointer.To(circuitPeerPeeringId.ID()),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.Properties.AuthorizationKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("address_prefix_ipv6"); ok {
		circuitId := expressroutecircuits.NewExpressRouteCircuitID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroupName, circuitPeeringId.CircuitName)

		circuit, err := circuitClient.Get(ctx, circuitId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", circuitId, err)
		}

		if circuit.Model != nil && circuit.Model.Properties != nil && circuit.Model.Properties.ExpressRoutePort != nil {
			return fmt.Errorf("`address_prefix_ipv6` cannot be set when ExpressRoute Circuit Connection with ExpressRoute Circuit based on ExpressRoute Port")
		} else {
			expressRouteCircuitConnectionParameters.Properties.IPv6CircuitConnectionConfig = &expressroutecircuitconnections.IPv6CircuitConnectionConfig{
				AddressPrefix: pointer.To(v.(string)),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, expressRouteCircuitConnectionParameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitConnectionRead(d, meta)
}

func resourceExpressRouteCircuitConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnections
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuitconnections.ParsePeeringConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("peering_id", commonids.NewExpressRouteCircuitPeeringID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteCircuitName, id.PeeringName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			d.Set("address_prefix_ipv4", props.AddressPrefix)

			// The ExpressRoute Circuit Connection API returns "*****************" for AuthorizationKey when it's changed from a valid value to `nil`
			// See more details from https://github.com/Azure/azure-rest-api-specs/issues/15030
			authorizationKey := ""
			if props.AuthorizationKey != nil && *props.AuthorizationKey != "*****************" {
				authorizationKey = *props.AuthorizationKey
			}
			d.Set("authorization_key", authorizationKey)

			addressPrefixIPv6 := ""
			if props.IPv6CircuitConnectionConfig != nil && props.IPv6CircuitConnectionConfig.AddressPrefix != nil {
				addressPrefixIPv6 = *props.IPv6CircuitConnectionConfig.AddressPrefix
			}
			d.Set("address_prefix_ipv6", addressPrefixIPv6)

			if props.PeerExpressRouteCircuitPeering != nil && props.PeerExpressRouteCircuitPeering.Id != nil {
				circuitPeerPeeringId, err := commonids.ParseExpressRouteCircuitPeeringIDInsensitively(*props.PeerExpressRouteCircuitPeering.Id)
				if err != nil {
					return err
				}
				d.Set("peer_peering_id", circuitPeerPeeringId.ID())
			}
		}
	}

	return nil
}

func resourceExpressRouteCircuitConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnections
	circuitClient := meta.(*clients.Client).Network.ExpressRouteCircuits
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuitconnections.ParsePeeringConnectionID(d.Id())
	if err != nil {
		return err
	}

	circuitPeeringId, err := commonids.ParseExpressRouteCircuitPeeringID(d.Get("peering_id").(string))
	if err != nil {
		return err
	}

	circuitPeerPeeringId, err := commonids.ParseExpressRouteCircuitPeeringID(d.Get("peer_peering_id").(string))
	if err != nil {
		return err
	}

	expressRouteCircuitConnectionParameters := expressroutecircuitconnections.ExpressRouteCircuitConnection{
		Name: pointer.To(id.ConnectionName),
		Properties: &expressroutecircuitconnections.ExpressRouteCircuitConnectionPropertiesFormat{
			AddressPrefix: pointer.To(d.Get("address_prefix_ipv4").(string)),
			ExpressRouteCircuitPeering: &expressroutecircuitconnections.SubResource{
				Id: pointer.To(circuitPeeringId.ID()),
			},
			PeerExpressRouteCircuitPeering: &expressroutecircuitconnections.SubResource{
				Id: pointer.To(circuitPeerPeeringId.ID()),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		expressRouteCircuitConnectionParameters.Properties.AuthorizationKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("address_prefix_ipv6"); ok {
		circuitId := expressroutecircuits.NewExpressRouteCircuitID(circuitPeeringId.SubscriptionId, circuitPeeringId.ResourceGroupName, circuitPeeringId.CircuitName)

		circuit, err := circuitClient.Get(ctx, circuitId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", circuitId, err)
		}

		if circuit.Model != nil && circuit.Model.Properties != nil && circuit.Model.Properties.ExpressRoutePort != nil {
			return fmt.Errorf("`address_prefix_ipv6` cannot be set when ExpressRoute Circuit Connection with ExpressRoute Circuit based on ExpressRoute Port")
		} else {
			expressRouteCircuitConnectionParameters.Properties.IPv6CircuitConnectionConfig = &expressroutecircuitconnections.IPv6CircuitConnectionConfig{
				AddressPrefix: pointer.To(v.(string)),
			}
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, expressRouteCircuitConnectionParameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceExpressRouteCircuitConnectionRead(d, meta)
}

func resourceExpressRouteCircuitConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnections
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuitconnections.ParsePeeringConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
