// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualHubConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubConnectionCreateOrUpdate,
		Read:   resourceVirtualHubConnectionRead,
		Update: resourceVirtualHubConnectionCreateOrUpdate,
		Delete: resourceVirtualHubConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseHubVirtualNetworkConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VirtualHubConnectionName,
			},

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"remote_virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
			},

			"internet_security_enabled": {
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
						"associated_route_table_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: virtualwans.ValidateHubRouteTableID,
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table", "routing.0.static_vnet_route"},
						},

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
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table", "routing.0.static_vnet_route"},
						},

						"static_vnet_local_route_override_criteria": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(virtualwans.VnetLocalRouteOverrideCriteriaContains),
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualwans.VnetLocalRouteOverrideCriteriaContains),
								string(virtualwans.VnetLocalRouteOverrideCriteriaEqual),
							}, false),
						},

						//lintignore:XS003
						"static_vnet_route": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"address_prefixes": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.IsCIDR,
										},
									},

									"next_hop_ip_address": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
								},
							},
							AtLeastOneOf: []string{"routing.0.associated_route_table_id", "routing.0.propagated_route_table", "routing.0.static_vnet_route"},
						},
					},
				},
			},
		},
	}
}

func resourceVirtualHubConnectionCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	id := virtualwans.NewHubVirtualNetworkConnectionID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, d.Get("name").(string))

	locks.ByName(virtualHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubId.VirtualHubName, virtualHubResourceName)

	remoteVirtualNetworkId, err := commonids.ParseVirtualNetworkID(d.Get("remote_virtual_network_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(remoteVirtualNetworkId.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(remoteVirtualNetworkId.VirtualNetworkName, VirtualNetworkResourceName)

	if d.IsNewResource() {
		existing, err := client.HubVirtualNetworkConnectionsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_hub_connection", id.ID())
		}
	}

	connection := virtualwans.HubVirtualNetworkConnection{
		Name: pointer.To(id.HubVirtualNetworkConnectionName),
		Properties: &virtualwans.HubVirtualNetworkConnectionProperties{
			RemoteVirtualNetwork: &virtualwans.SubResource{
				Id: pointer.To(remoteVirtualNetworkId.ID()),
			},
			EnableInternetSecurity: pointer.To(d.Get("internet_security_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("routing"); ok {
		connection.Properties.RoutingConfiguration = expandVirtualHubConnectionRouting(v.([]interface{}))
	}

	if err := client.HubVirtualNetworkConnectionsCreateOrUpdateThenPoll(ctx, id, connection); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualwans.ProvisioningStateUpdating)},
		Target:     []string{string(virtualwans.ProvisioningStateSucceeded)},
		Refresh:    virtualHubConnectionProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubConnectionRead(d, meta)
}

func resourceVirtualHubConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubVirtualNetworkConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.HubVirtualNetworkConnectionName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("internet_security_enabled", props.EnableInternetSecurity)
			remoteVirtualNetworkId := ""
			if props.RemoteVirtualNetwork != nil && props.RemoteVirtualNetwork.Id != nil {
				remoteVirtualNetworkId = *props.RemoteVirtualNetwork.Id
			}
			d.Set("remote_virtual_network_id", remoteVirtualNetworkId)

			if err := d.Set("routing", flattenVirtualHubConnectionRouting(props.RoutingConfiguration)); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
		}
	}

	return nil
}

func resourceVirtualHubConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	if err := client.HubVirtualNetworkConnectionsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVirtualHubConnectionRouting(input []interface{}) *virtualwans.RoutingConfiguration {
	if len(input) == 0 {
		return &virtualwans.RoutingConfiguration{}
	}

	v := input[0].(map[string]interface{})

	result := &virtualwans.RoutingConfiguration{
		VnetRoutes: &virtualwans.VnetRoute{
			StaticRoutes: expandVirtualHubConnectionVnetStaticRoute(v["static_vnet_route"].([]interface{})),
			StaticRoutesConfig: &virtualwans.StaticRoutesConfig{
				VnetLocalRouteOverrideCriteria: pointer.To(virtualwans.VnetLocalRouteOverrideCriteria(v["static_vnet_local_route_override_criteria"].(string))),
			},
		},
	}

	if associatedRouteTableId := v["associated_route_table_id"].(string); associatedRouteTableId != "" {
		result.AssociatedRouteTable = &virtualwans.SubResource{
			Id: pointer.To(associatedRouteTableId),
		}
	}

	if inboundRouteMapId := v["inbound_route_map_id"].(string); inboundRouteMapId != "" {
		result.InboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(inboundRouteMapId),
		}
	}

	if outboundRouteMapId := v["outbound_route_map_id"].(string); outboundRouteMapId != "" {
		result.OutboundRouteMap = &virtualwans.SubResource{
			Id: pointer.To(outboundRouteMapId),
		}
	}

	if propagatedRouteTable := v["propagated_route_table"].([]interface{}); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandVirtualHubConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return result
}

func expandVirtualHubConnectionPropagatedRouteTable(input []interface{}) *virtualwans.PropagatedRouteTable {
	if len(input) == 0 {
		return &virtualwans.PropagatedRouteTable{}
	}

	v := input[0].(map[string]interface{})

	result := virtualwans.PropagatedRouteTable{}

	if labels := v["labels"].(*pluginsdk.Set).List(); len(labels) != 0 {
		result.Labels = utils.ExpandStringSlice(labels)
	}

	if routeTableIds := v["route_table_ids"].([]interface{}); len(routeTableIds) != 0 {
		result.Ids = expandIDsToVirtualWANSubResources(routeTableIds)
	}

	return &result
}

func expandVirtualHubConnectionVnetStaticRoute(input []interface{}) *[]virtualwans.StaticRoute {
	if len(input) == 0 {
		return nil
	}

	results := make([]virtualwans.StaticRoute, 0)

	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})

		result := virtualwans.StaticRoute{}

		if name := v["name"].(string); name != "" {
			result.Name = pointer.To(name)
		}

		if addressPrefixes := v["address_prefixes"].(*pluginsdk.Set).List(); len(addressPrefixes) != 0 {
			result.AddressPrefixes = utils.ExpandStringSlice(addressPrefixes)
		}

		if nextHopIPAddress := v["next_hop_ip_address"].(string); nextHopIPAddress != "" {
			result.NextHopIPAddress = pointer.To(nextHopIPAddress)
		}

		results = append(results, result)
	}

	return &results
}

func expandIDsToVirtualWANSubResources(input []interface{}) *[]virtualwans.SubResource {
	ids := make([]virtualwans.SubResource, 0)

	for _, v := range input {
		ids = append(ids, virtualwans.SubResource{
			Id: pointer.To(v.(string)),
		})
	}

	return &ids
}

func flattenVirtualHubConnectionRouting(input *virtualwans.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	associatedRouteTableId := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.Id != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.Id
	}

	inboundRouteMapId := ""
	if input.InboundRouteMap != nil && input.InboundRouteMap.Id != nil {
		inboundRouteMapId = *input.InboundRouteMap.Id
	}

	outboundRouteMapId := ""
	if input.OutboundRouteMap != nil && input.OutboundRouteMap.Id != nil {
		outboundRouteMapId = *input.OutboundRouteMap.Id
	}

	staticVnetLocalRouteOverrideCriteria := ""
	if input.VnetRoutes != nil && input.VnetRoutes.StaticRoutesConfig != nil && input.VnetRoutes.StaticRoutesConfig.VnetLocalRouteOverrideCriteria != nil {
		staticVnetLocalRouteOverrideCriteria = string(*input.VnetRoutes.StaticRoutesConfig.VnetLocalRouteOverrideCriteria)
	}

	return []interface{}{
		map[string]interface{}{
			"associated_route_table_id":                 associatedRouteTableId,
			"inbound_route_map_id":                      inboundRouteMapId,
			"outbound_route_map_id":                     outboundRouteMapId,
			"propagated_route_table":                    flattenVirtualHubConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"static_vnet_route":                         flattenVirtualHubConnectionVnetStaticRoute(input.VnetRoutes),
			"static_vnet_local_route_override_criteria": staticVnetLocalRouteOverrideCriteria,
		},
	}
}

func flattenVirtualHubConnectionPropagatedRouteTable(input *virtualwans.PropagatedRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	labels := make([]interface{}, 0)
	if input.Labels != nil {
		labels = utils.FlattenStringSlice(input.Labels)
	}

	routeTableIds := make([]interface{}, 0)
	if input.Ids != nil {
		routeTableIds = flattenVirtualWANSubResourcesToIDs(input.Ids)
	}

	return []interface{}{
		map[string]interface{}{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func flattenVirtualHubConnectionVnetStaticRoute(input *virtualwans.VnetRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.StaticRoutes == nil {
		return results
	}

	for _, item := range *input.StaticRoutes {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var nextHopIpAddress string
		if item.NextHopIPAddress != nil {
			nextHopIpAddress = *item.NextHopIPAddress
		}

		addressPrefixes := make([]interface{}, 0)
		if item.AddressPrefixes != nil {
			addressPrefixes = utils.FlattenStringSlice(item.AddressPrefixes)
		}

		v := map[string]interface{}{
			"name":                name,
			"address_prefixes":    addressPrefixes,
			"next_hop_ip_address": nextHopIpAddress,
		}

		results = append(results, v)
	}

	return results
}

func flattenVirtualWANSubResourcesToIDs(input *[]virtualwans.SubResource) []interface{} {
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

func virtualHubConnectionProvisioningStateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id virtualwans.HubVirtualNetworkConnectionId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.HubVirtualNetworkConnectionsGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
		}

		if res.Model == nil {
			return res, "", fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		if res.Model.Properties == nil {
			return res, "", fmt.Errorf("retrieving %s: `properties` was nil", id)
		}
		if res.Model.Properties.ProvisioningState == nil {
			return res, "", fmt.Errorf("retrieving %s: `properties.provisioningState` was nil", id)
		}
		return res, string(*res.Model.Properties.ProvisioningState), nil
	}
}
