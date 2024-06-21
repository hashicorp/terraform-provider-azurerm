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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const virtualHubResourceName = "azurerm_virtual_hub"

func resourceVirtualHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubCreate,
		Read:   resourceVirtualHubRead,
		Update: resourceVirtualHubUpdate,
		Delete: resourceVirtualHubDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseVirtualHubID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.VirtualHubName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"address_prefix": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.CIDR,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Basic",
					"Standard",
				}, false),
			},

			"virtual_wan_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualWANID,
			},

			"virtual_router_asn": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"virtual_router_ips": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"route": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"address_prefixes": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},
						"next_hop_ip_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.IPv4Address,
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"hub_routing_preference": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(virtualwans.HubRoutingPreferenceExpressRoute),
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualwans.HubRoutingPreferenceExpressRoute),
					string(virtualwans.HubRoutingPreferenceVpnGateway),
					string(virtualwans.HubRoutingPreferenceASPath),
				}, false),
			},

			"default_route_table_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"virtual_router_auto_scale_min_capacity": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(2),
				Default:      2,
			},
		},
	}
}

func resourceVirtualHubCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("deadline is not properly set for Virtual Hub")
	}

	id := virtualwans.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	existing, err := client.VirtualHubsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_hub", id.ID())
	}

	parameters := virtualwans.VirtualHub{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VirtualHubProperties{
			RouteTable:           expandVirtualHubRoute(d.Get("route").(*pluginsdk.Set).List()),
			HubRoutingPreference: pointer.To(virtualwans.HubRoutingPreference(d.Get("hub_routing_preference").(string))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("address_prefix"); ok {
		parameters.Properties.AddressPrefix = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("sku"); ok {
		parameters.Properties.Sku = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("virtual_wan_id"); ok {
		parameters.Properties.VirtualWAN = &virtualwans.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("virtual_router_auto_scale_min_capacity"); ok {
		parameters.Properties.VirtualRouterAutoScaleConfiguration = &virtualwans.VirtualRouterAutoScaleConfiguration{
			MinCapacity: pointer.To(int64(v.(int))),
		}
	}

	if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Hub returns provisioned while the routing state is still "provisioning". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, client, id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s provisioning route: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRead(d, meta)
}

func resourceVirtualHubUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("deadline is not properly set for Virtual Hub")
	}

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	existing, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("route") {
		payload.Properties.RouteTable = expandVirtualHubRoute(d.Get("route").(*pluginsdk.Set).List())
	}

	if d.HasChange("hub_routing_preference") {
		payload.Properties.HubRoutingPreference = pointer.To(virtualwans.HubRoutingPreference(d.Get("hub_routing_preference").(string)))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("virtual_router_auto_scale_min_capacity") {
		if v, ok := d.GetOk("virtual_router_auto_scale_min_capacity"); ok {
			payload.Properties.VirtualRouterAutoScaleConfiguration = &virtualwans.VirtualRouterAutoScaleConfiguration{
				MinCapacity: pointer.To(int64(v.(int))),
			}
		}
	}

	if err := client.VirtualHubsCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Hub returns provisioned while the routing state is still "provisioning". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, client, *id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s provisioning route: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRead(d, meta)
}

func resourceVirtualHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualHubsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualHubName)
	d.Set("resource_group_name", id.ResourceGroupName)

	defaultRouteTable := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, "defaultRouteTable")
	d.Set("default_route_table_id", defaultRouteTable.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {

			d.Set("address_prefix", props.AddressPrefix)
			d.Set("sku", props.Sku)

			if err := d.Set("route", flattenVirtualHubRoute(props.RouteTable)); err != nil {
				return fmt.Errorf("setting `route`: %+v", err)
			}

			d.Set("hub_routing_preference", string(pointer.From(props.HubRoutingPreference)))

			var virtualWanId *string
			if props.VirtualWAN != nil {
				virtualWanId = props.VirtualWAN.Id
			}
			d.Set("virtual_wan_id", virtualWanId)

			var virtualRouterAsn *int64
			if props.VirtualRouterAsn != nil {
				virtualRouterAsn = props.VirtualRouterAsn
			}
			d.Set("virtual_router_asn", virtualRouterAsn)

			var virtualRouterIps *[]string
			if props.VirtualRouterIPs != nil {
				virtualRouterIps = props.VirtualRouterIPs
			}
			d.Set("virtual_router_ips", virtualRouterIps)

			d.Set("virtual_router_auto_scale_min_capacity", props.VirtualRouterAutoScaleConfiguration.MinCapacity)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceVirtualHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	if err := client.VirtualHubsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVirtualHubRoute(input []interface{}) *virtualwans.VirtualHubRouteTable {
	if len(input) == 0 {
		return nil
	}

	results := make([]virtualwans.VirtualHubRoute, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})
		addressPrefixes := v["address_prefixes"].([]interface{})
		nextHopIpAddress := v["next_hop_ip_address"].(string)

		results = append(results, virtualwans.VirtualHubRoute{
			AddressPrefixes:  utils.ExpandStringSlice(addressPrefixes),
			NextHopIPAddress: pointer.To(nextHopIpAddress),
		})
	}

	result := virtualwans.VirtualHubRouteTable{
		Routes: &results,
	}

	return &result
}

func flattenVirtualHubRoute(input *virtualwans.VirtualHubRouteTable) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Routes == nil {
		return results
	}

	for _, item := range *input.Routes {
		addressPrefixes := utils.FlattenStringSlice(item.AddressPrefixes)
		nextHopIpAddress := ""

		if item.NextHopIPAddress != nil {
			nextHopIpAddress = *item.NextHopIPAddress
		}

		results = append(results, map[string]interface{}{
			"address_prefixes":    addressPrefixes,
			"next_hop_ip_address": nextHopIpAddress,
		})
	}

	return results
}

func virtualHubCreateRefreshFunc(ctx context.Context, client *virtualwans.VirtualWANsClient, id virtualwans.VirtualHubId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.VirtualHubsGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("%s doesn't exist", id)
			}
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Model == nil {
			return nil, "", fmt.Errorf("retrieving %s: `model` was nil", id)
		}
		if res.Model.Properties == nil {
			return nil, "", fmt.Errorf("retrieving %s: `properties` was nil", id)
		}

		state := res.Model.Properties.RoutingState
		if state != nil && *state == "Failed" {
			return nil, "", fmt.Errorf("failed to provision routing on %s", id)
		}
		return res, string(*res.Model.Properties.RoutingState), nil
	}
}
