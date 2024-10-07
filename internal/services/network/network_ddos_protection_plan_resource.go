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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const ddosProtectionPlanResourceName = "azurerm_network_ddos_protection_plan"

func resourceNetworkDDoSProtectionPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkDDoSProtectionPlanCreate,
		Read:   resourceNetworkDDoSProtectionPlanRead,
		Update: resourceNetworkDDoSProtectionPlanUpdate,
		Delete: resourceNetworkDDoSProtectionPlanDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := ddosprotectionplans.ParseDdosProtectionPlanID(id)
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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"virtual_network_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceNetworkDDoSProtectionPlanCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetsToLock, err := expandNetworkDDoSProtectionPlanVnetNames(d.Get("virtual_network_ids").([]interface{}))
	if err != nil {
		return fmt.Errorf("extracting names of Virtual Network: %+v", err)
	}

	id := ddosprotectionplans.NewDdosProtectionPlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.DdosProtectionPlanName, ddosProtectionPlanResourceName)
	defer locks.UnlockByName(id.DdosProtectionPlanName, ddosProtectionPlanResourceName)
	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_ddos_protection_plan", id.ID())
	}

	payload := ddosprotectionplans.DdosProtectionPlan{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkDDoSProtectionPlanRead(d, meta)
}

func resourceNetworkDDoSProtectionPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetsToLock, err := expandNetworkDDoSProtectionPlanVnetNames(d.Get("virtual_network_ids").([]interface{}))
	if err != nil {
		return fmt.Errorf("extracting names of Virtual Network: %+v", err)
	}

	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(d.Id())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	locks.ByName(id.DdosProtectionPlanName, ddosProtectionPlanResourceName)
	defer locks.UnlockByName(id.DdosProtectionPlanName, ddosProtectionPlanResourceName)
	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	existing, err := client.Get(ctx, *id)
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

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkDDoSProtectionPlanRead(d, meta)
}

func resourceNetworkDDoSProtectionPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.DdosProtectionPlanName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			vNetIDs := flattenNetworkDDoSProtectionPlanVirtualNetworkIDs(props.VirtualNetworks)
			if err := d.Set("virtual_network_ids", vNetIDs); err != nil {
				return fmt.Errorf("setting `virtual_network_ids`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceNetworkDDoSProtectionPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ddosprotectionplans.ParseDdosProtectionPlanID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}
	// if there's no VirtualNetworks configured, it's possible for this to be nil
	subResources := existing.Model.Properties.VirtualNetworks
	virtualNetworksNamesToLock, err := extractVnetNames(subResources)
	if err != nil {
		return fmt.Errorf("extracting names of Virtual Network: %+v", err)
	}

	locks.ByName(id.DdosProtectionPlanName, ddosProtectionPlanResourceName)
	defer locks.UnlockByName(id.DdosProtectionPlanName, ddosProtectionPlanResourceName)

	locks.MultipleByName(virtualNetworksNamesToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(virtualNetworksNamesToLock, VirtualNetworkResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetworkDDoSProtectionPlanVnetNames(input []interface{}) (*[]string, error) {
	vnetNames := make([]string, 0)

	for _, vnetID := range input {
		vnetResourceID, err := commonids.ParseVirtualNetworkID(vnetID.(string))
		if err != nil {
			return nil, err
		}

		if !utils.SliceContainsValue(vnetNames, vnetResourceID.VirtualNetworkName) {
			vnetNames = append(vnetNames, vnetResourceID.VirtualNetworkName)
		}
	}

	return &vnetNames, nil
}

func flattenNetworkDDoSProtectionPlanVirtualNetworkIDs(input *[]ddosprotectionplans.SubResource) []string {
	vnetIDs := make([]string, 0)
	if input == nil {
		return vnetIDs
	}

	// if-continue is used to simplify the deeply nested if-else statement.
	for _, subRes := range *input {
		if subRes.Id != nil {
			vnetIDs = append(vnetIDs, *subRes.Id)
		}
	}

	return vnetIDs
}

func extractVnetNames(input *[]ddosprotectionplans.SubResource) (*[]string, error) {
	vnetNames := make([]string, 0)

	if input != nil {
		for _, subresource := range *input {
			if subresource.Id == nil {
				continue
			}

			id, err := commonids.ParseVirtualNetworkIDInsensitively(*subresource.Id)
			if err != nil {
				return nil, err
			}

			if !utils.SliceContainsValue(vnetNames, id.VirtualNetworkName) {
				vnetNames = append(vnetNames, id.VirtualNetworkName)
			}
		}
	}

	return &vnetNames, nil
}
