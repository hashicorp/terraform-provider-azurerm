// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutegateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceExpressRouteConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteConnectionCreate,
		Read:   resourceExpressRouteConnectionRead,
		Update: resourceExpressRouteConnectionUpdate,
		Delete: resourceExpressRouteConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressrouteconnections.ParseExpressRouteConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteConnectionName,
			},

			"express_route_circuit_peering_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateExpressRouteCircuitPeeringID,
			},

			"express_route_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: expressroutegateways.ValidateExpressRouteGatewayID,
			},

			"authorization_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_internet_security": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"private_link_fast_path_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"express_route_gateway_bypass_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"routing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"inbound_route_map_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: virtualwans.ValidateRouteMapID,
						},

						"outbound_route_map_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: virtualwans.ValidateRouteMapID,
						},

						"associated_route_table_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: virtualwans.ValidateHubRouteTableID,
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table"},
						},

						"propagated_route_table": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"labels": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
										AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
									},

									"route_table_ids": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: virtualwans.ValidateHubRouteTableID,
										},
										AtLeastOneOf: []string{"routing.0.propagated_route_table.0.labels", "routing.0.propagated_route_table.0.route_table_ids"},
									},
								},
							},
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table"},
						},
					},
				},
			},

			"routing_weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 32000),
			},
		},
	}
}

func resourceExpressRouteConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnections
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	expressRouteGatewayId, err := expressroutegateways.ParseExpressRouteGatewayID(d.Get("express_route_gateway_id").(string))
	if err != nil {
		return err
	}

	id := expressrouteconnections.NewExpressRouteConnectionID(expressRouteGatewayId.SubscriptionId, expressRouteGatewayId.ResourceGroupName, expressRouteGatewayId.ExpressRouteGatewayName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_express_route_connection", id.ID())
	}

	parameters := expressrouteconnections.ExpressRouteConnection{
		Name: id.ExpressRouteConnectionName,
		Properties: &expressrouteconnections.ExpressRouteConnectionProperties{
			ExpressRouteCircuitPeering: expressrouteconnections.ExpressRouteCircuitPeeringId{
				Id: pointer.To(d.Get("express_route_circuit_peering_id").(string)),
			},
			EnableInternetSecurity:    pointer.To(d.Get("enable_internet_security").(bool)),
			RoutingConfiguration:      expandExpressRouteConnectionRouting(d.Get("routing").([]interface{})),
			RoutingWeight:             pointer.To(int64(d.Get("routing_weight").(int))),
			ExpressRouteGatewayBypass: pointer.To(d.Get("express_route_gateway_bypass_enabled").(bool)),
		},
	}

	privateLinkFastPath := d.GetRawConfig().AsValueMap()["private_link_fast_path_enabled"]
	if !privateLinkFastPath.IsNull() {
		if d.Get("private_link_fast_path_enabled").(bool) && !d.Get("express_route_gateway_bypass_enabled").(bool) {
			return fmt.Errorf("`express_route_gateway_bypass_enabled` must be enabled when `private_link_fast_path_enabled` is set to `true`")
		}
		parameters.Properties.EnablePrivateLinkFastPath = pointer.To(d.Get("private_link_fast_path_enabled").(bool))
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		parameters.Properties.AuthorizationKey = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteConnectionRead(d, meta)
}

func resourceExpressRouteConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnections
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteconnections.ParseExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ExpressRouteConnectionName)
	d.Set("express_route_gateway_id", expressroutegateways.NewExpressRouteGatewayID(id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteGatewayName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("routing_weight", props.RoutingWeight)
			d.Set("authorization_key", props.AuthorizationKey)
			d.Set("enable_internet_security", props.EnableInternetSecurity)
			d.Set("private_link_fast_path_enabled", pointer.From(props.EnablePrivateLinkFastPath))

			if props.ExpressRouteGatewayBypass != nil {
				d.Set("express_route_gateway_bypass_enabled", props.ExpressRouteGatewayBypass)
			}

			circuitPeeringID := ""
			if v := props.ExpressRouteCircuitPeering.Id; v != nil {
				circuitPeeringID = *v
			}
			peeringId, err := commonids.ParseExpressRouteCircuitPeeringIDInsensitively(circuitPeeringID)
			if err != nil {
				return err
			}
			d.Set("express_route_circuit_peering_id", peeringId.ID())

			routing, err := flattenExpressRouteConnectionRouting(props.RoutingConfiguration)
			if err != nil {
				return err
			}
			if err := d.Set("routing", routing); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
		}
	}

	return nil
}

func resourceExpressRouteConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnections
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteconnections.ParseExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	parameters := expressrouteconnections.ExpressRouteConnection{
		Name: id.ExpressRouteConnectionName,
		Properties: &expressrouteconnections.ExpressRouteConnectionProperties{
			ExpressRouteCircuitPeering: expressrouteconnections.ExpressRouteCircuitPeeringId{
				Id: pointer.To(d.Get("express_route_circuit_peering_id").(string)),
			},
			EnableInternetSecurity:    pointer.To(d.Get("enable_internet_security").(bool)),
			RoutingConfiguration:      expandExpressRouteConnectionRouting(d.Get("routing").([]interface{})),
			RoutingWeight:             pointer.To(int64(d.Get("routing_weight").(int))),
			ExpressRouteGatewayBypass: pointer.To(d.Get("express_route_gateway_bypass_enabled").(bool)),
		},
	}

	privateLinkFastPath := d.GetRawConfig().AsValueMap()["private_link_fast_path_enabled"]
	if !privateLinkFastPath.IsNull() {
		if d.Get("private_link_fast_path_enabled").(bool) && !d.Get("express_route_gateway_bypass_enabled").(bool) {
			return fmt.Errorf("`express_route_gateway_bypass_enabled` must be enabled when `private_link_fast_path_enabled` is set to `true`")
		}
		parameters.Properties.EnablePrivateLinkFastPath = pointer.To(d.Get("private_link_fast_path_enabled").(bool))
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		parameters.Properties.AuthorizationKey = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceExpressRouteConnectionRead(d, meta)
}

func resourceExpressRouteConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnections
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteconnections.ParseExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandExpressRouteConnectionRouting(input []interface{}) *expressrouteconnections.RoutingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return &expressrouteconnections.RoutingConfiguration{}
	}

	v := input[0].(map[string]interface{})
	result := expressrouteconnections.RoutingConfiguration{}

	if associatedRouteTableId := v["associated_route_table_id"].(string); associatedRouteTableId != "" {
		result.AssociatedRouteTable = &expressrouteconnections.SubResource{
			Id: pointer.To(associatedRouteTableId),
		}
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		result.InboundRouteMap = &expressrouteconnections.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		result.OutboundRouteMap = &expressrouteconnections.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	if propagatedRouteTable := v["propagated_route_table"].([]interface{}); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandExpressRouteConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return &result
}

func expandExpressRouteConnectionPropagatedRouteTable(input []interface{}) *expressrouteconnections.PropagatedRouteTable {
	if len(input) == 0 || input[0] == nil {
		return &expressrouteconnections.PropagatedRouteTable{}
	}

	v := input[0].(map[string]interface{})

	result := expressrouteconnections.PropagatedRouteTable{}

	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = utils.ExpandStringSlice(labels)
	}

	if routeTableIds := v["route_table_ids"].([]interface{}); len(routeTableIds) != 0 {
		result.Ids = expandExpressRouteIDsToSubResources(routeTableIds)
	}

	return &result
}

func expandExpressRouteIDsToSubResources(input []interface{}) *[]expressrouteconnections.SubResource {
	ids := make([]expressrouteconnections.SubResource, 0)

	for _, v := range input {
		ids = append(ids, expressrouteconnections.SubResource{
			Id: pointer.To(v.(string)),
		})
	}

	return &ids
}

func flattenExpressRouteConnectionRouting(input *expressrouteconnections.RoutingConfiguration) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	associatedRouteTableId := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.Id
	}
	routeTableId, err := virtualwans.ParseHubRouteTableIDInsensitively(associatedRouteTableId)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"associated_route_table_id": routeTableId.ID(),
		"propagated_route_table":    flattenExpressRouteConnectionPropagatedRouteTable(input.PropagatedRouteTables),
	}

	if input.InboundRouteMap != nil && input.InboundRouteMap.Id != nil {
		result["inbound_route_map_id"] = input.InboundRouteMap.Id
	}

	if input.OutboundRouteMap != nil && input.OutboundRouteMap.Id != nil {
		result["outbound_route_map_id"] = input.OutboundRouteMap.Id
	}

	return []interface{}{result}, nil
}

func flattenExpressRouteConnectionPropagatedRouteTable(input *expressrouteconnections.PropagatedRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	labels := make([]interface{}, 0)
	if input.Labels != nil {
		labels = utils.FlattenStringSlice(input.Labels)
	}

	routeTableIds := make([]interface{}, 0)
	if input.Ids != nil {
		routeTableIds = flattenExpressRouteSubResourcesToIDs(input.Ids)
	}

	return []interface{}{
		map[string]interface{}{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func flattenExpressRouteSubResourcesToIDs(input *[]expressrouteconnections.SubResource) []interface{} {
	ids := make([]interface{}, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}
