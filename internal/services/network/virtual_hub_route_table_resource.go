// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualHubRouteTable() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubRouteTableCreate,
		Read:   resourceVirtualHubRouteTableRead,
		Update: resourceVirtualHubRouteTableUpdate,
		Delete: resourceVirtualHubRouteTableDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualwans.ParseHubRouteTableID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.HubRouteTableName,
			},

			"virtual_hub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"labels": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"route": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"destinations": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"destinations_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CIDR",
								"ResourceId",
								"Service",
							}, false),
						},

						"next_hop": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"next_hop_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "ResourceId",
							ValidateFunc: validation.StringInSlice([]string{
								"ResourceId",
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceVirtualHubRouteTableCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id := virtualwans.NewHubRouteTableID(virtHubId.SubscriptionId, virtHubId.ResourceGroupName, virtHubId.VirtualHubName, d.Get("name").(string))

	existing, err := client.HubRouteTablesGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_hub_route_table", id.ID())
	}

	parameters := virtualwans.HubRouteTable{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.HubRouteTableProperties{
			Labels: utils.ExpandStringSlice(d.Get("labels").(*pluginsdk.Set).List()),
			Routes: expandVirtualHubRouteTableHubRoutes(d.Get("route").(*pluginsdk.Set).List()),
		},
	}

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRead(d, meta)
}

func resourceVirtualHubRouteTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id, err := virtualwans.ParseHubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.HubRouteTablesGet(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
	}

	payload := existing.Model

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("labels") {
		payload.Properties.Labels = utils.ExpandStringSlice(d.Get("labels").(*pluginsdk.Set).List())
	}

	if d.HasChange("route") {
		payload.Properties.Routes = expandVirtualHubRouteTableHubRoutes(d.Get("route").(*pluginsdk.Set).List())
	}

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRead(d, meta)
}

func resourceVirtualHubRouteTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubRouteTablesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.HubRouteTableName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("labels", utils.FlattenStringSlice(props.Labels))

			if err := d.Set("route", flattenVirtualHubRouteTableHubRoutes(props.Routes)); err != nil {
				return fmt.Errorf("setting `route`: %+v", err)
			}
		}
	}

	return nil
}

func resourceVirtualHubRouteTableDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubRouteTableID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	if err := client.HubRouteTablesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandVirtualHubRouteTableHubRoutes(input []interface{}) *[]virtualwans.HubRoute {
	results := make([]virtualwans.HubRoute, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := virtualwans.HubRoute{
			Name:            v["name"].(string),
			DestinationType: v["destinations_type"].(string),
			Destinations:    pointer.From(utils.ExpandStringSlice(v["destinations"].(*pluginsdk.Set).List())),
			NextHopType:     v["next_hop_type"].(string),
			NextHop:         v["next_hop"].(string),
		}

		results = append(results, result)
	}

	return &results
}

func flattenVirtualHubRouteTableHubRoutes(input *[]virtualwans.HubRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := map[string]interface{}{
			"name":              item.Name,
			"destinations":      utils.FlattenStringSlice(&item.Destinations),
			"destinations_type": item.DestinationType,
			"next_hop":          item.NextHop,
			"next_hop_type":     item.NextHopType,
		}

		results = append(results, v)
	}

	return results
}
