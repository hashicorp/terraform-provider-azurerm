// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routetables"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var routeTableResourceName = "azurerm_route_table"

func resourceRouteTable() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceRouteTableCreate,
		Read:   resourceRouteTableRead,
		Update: resourceRouteTableUpdate,
		Delete: resourceRouteTableDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := routetables.ParseRouteTableID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RouteTableName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"route": {
				Type:       pluginsdk.TypeSet,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				Computed:   true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.RouteName,
						},

						"address_prefix": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"next_hop_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(routetables.RouteNextHopTypeVirtualNetworkGateway),
								string(routetables.RouteNextHopTypeVnetLocal),
								string(routetables.RouteNextHopTypeInternet),
								string(routetables.RouteNextHopTypeVirtualAppliance),
								string(routetables.RouteNextHopTypeNone),
							}, false),
						},

						"next_hop_in_ip_address": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"bgp_route_propagation_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"subnets": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["bgp_route_propagation_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ConflictsWith: []string{
				"disable_bgp_route_propagation",
			},
		}
		resource.Schema["disable_bgp_route_propagation"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ConflictsWith: []string{
				"bgp_route_propagation_enabled",
			},
			Deprecated: "The property `disable_bgp_route_propagation` has been superseded by the property `bgp_route_propagation_enabled` and will be removed in v4.0 of the AzureRM Provider.",
		}
	}

	return resource
}

func resourceRouteTableCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteTables
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := routetables.NewRouteTableID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, routetables.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_route_table", id.ID())
	}

	bgpRoutePropagationEnabled := d.Get("bgp_route_propagation_enabled").(bool)

	if !features.FourPointOhBeta() {
		// need to set default back to true for 3.x, this triggers ineffassign linter, so ignoring for now
		bgpRoutePropagationEnabled = true // nolint: ineffassign
		bgpRoutePropagationEnabled = !d.Get("disable_bgp_route_propagation").(bool)
	}

	routeSet := routetables.RouteTable{
		Name:     &id.RouteTableName,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &routetables.RouteTablePropertiesFormat{
			Routes:                     expandRouteTableRoutes(d),
			DisableBgpRoutePropagation: pointer.To(!bgpRoutePropagationEnabled),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, routeSet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceRouteTableRead(d, meta)
}

func resourceRouteTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteTables
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routetables.ParseRouteTableID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, routetables.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("route") {
		payload.Properties.Routes = expandRouteTableRoutes(d)
	}

	if !features.FourPointOhBeta() {
		if d.HasChange("disable_bgp_route_propagation") {
			payload.Properties.DisableBgpRoutePropagation = pointer.To(d.Get("disable_bgp_route_propagation").(bool))
		}
	}

	if d.HasChange("bgp_route_propagation_enabled") {
		payload.Properties.DisableBgpRoutePropagation = pointer.To(!d.Get("bgp_route_propagation_enabled").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceRouteTableRead(d, meta)
}

func resourceRouteTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteTables
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routetables.ParseRouteTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, routetables.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RouteTableName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if !features.FourPointOhBeta() {
				d.Set("disable_bgp_route_propagation", props.DisableBgpRoutePropagation)
			}
			d.Set("bgp_route_propagation_enabled", !pointer.From(props.DisableBgpRoutePropagation))
			if err := d.Set("route", flattenRouteTableRoutes(props.Routes)); err != nil {
				return err
			}

			if err := d.Set("subnets", flattenRouteTableSubnets(props.Subnets)); err != nil {
				return err
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceRouteTableDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RouteTables
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routetables.ParseRouteTableID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandRouteTableRoutes(d *pluginsdk.ResourceData) *[]routetables.Route {
	configs := d.Get("route").(*pluginsdk.Set).List()
	routes := make([]routetables.Route, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		route := routetables.Route{
			Name: utils.String(data["name"].(string)),
			Properties: &routetables.RoutePropertiesFormat{
				AddressPrefix: utils.String(data["address_prefix"].(string)),
				NextHopType:   routetables.RouteNextHopType(data["next_hop_type"].(string)),
			},
		}

		if v := data["next_hop_in_ip_address"].(string); v != "" {
			route.Properties.NextHopIPAddress = &v
		}

		routes = append(routes, route)
	}

	return &routes
}

func flattenRouteTableRoutes(input *[]routetables.Route) []interface{} {
	results := make([]interface{}, 0)

	if routes := input; routes != nil {
		for _, route := range *routes {
			r := make(map[string]interface{})

			r["name"] = *route.Name

			if props := route.Properties; props != nil {
				r["address_prefix"] = *props.AddressPrefix
				r["next_hop_type"] = string(props.NextHopType)
				if ip := props.NextHopIPAddress; ip != nil && *ip != "" {
					r["next_hop_in_ip_address"] = *ip
				}
			}

			results = append(results, r)
		}
	}

	return results
}

func flattenRouteTableSubnets(subnets *[]routetables.Subnet) []string {
	output := make([]string, 0)

	if subnets != nil {
		for _, subnet := range *subnets {
			output = append(output, *subnet.Id)
		}
	}

	return output
}
