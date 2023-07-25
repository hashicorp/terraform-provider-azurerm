// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

const azureNetworkDDoSProtectionPlanResourceName = "azurerm_network_ddos_protection_plan"

func resourceNetworkDDoSProtectionPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkDDoSProtectionPlanCreateUpdate,
		Read:   resourceNetworkDDoSProtectionPlanRead,
		Update: resourceNetworkDDoSProtectionPlanCreateUpdate,
		Delete: resourceNetworkDDoSProtectionPlanDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DdosProtectionPlanID(id)
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

			"tags": tags.Schema(),
		},
	}
}

func resourceNetworkDDoSProtectionPlanCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DDOSProtectionPlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DDoS protection plan creation")

	id := parse.NewDdosProtectionPlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_network_ddos_protection_plan", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	vnetsToLock, err := expandNetworkDDoSProtectionPlanVnetNames(d)
	if err != nil {
		return fmt.Errorf("extracting names of Virtual Network: %+v", err)
	}

	locks.ByName(id.Name, azureNetworkDDoSProtectionPlanResourceName)
	defer locks.UnlockByName(id.Name, azureNetworkDDoSProtectionPlanResourceName)

	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	parameters := network.DdosProtectionPlan{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkDDoSProtectionPlanRead(d, meta)
}

func resourceNetworkDDoSProtectionPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DDOSProtectionPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DdosProtectionPlanID(d.Id())
	if err != nil {
		return err
	}

	plan, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(plan.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := plan.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := plan.DdosProtectionPlanPropertiesFormat; props != nil {
		vNetIDs := flattenNetworkDDoSProtectionPlanVirtualNetworkIDs(props.VirtualNetworks)
		if err := d.Set("virtual_network_ids", vNetIDs); err != nil {
			return fmt.Errorf("setting `virtual_network_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, plan.Tags)
}

func resourceNetworkDDoSProtectionPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.DDOSProtectionPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DdosProtectionPlanID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	vnetsToLock, err := extractVnetNames(d)
	if err != nil {
		return fmt.Errorf("extracting names of Virtual Network: %+v", err)
	}

	locks.ByName(id.Name, azureNetworkDDoSProtectionPlanResourceName)
	defer locks.UnlockByName(id.Name, azureNetworkDDoSProtectionPlanResourceName)

	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}

func expandNetworkDDoSProtectionPlanVnetNames(d *pluginsdk.ResourceData) (*[]string, error) {
	vnetIDs := d.Get("virtual_network_ids").([]interface{})
	vnetNames := make([]string, 0)

	for _, vnetID := range vnetIDs {
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

func flattenNetworkDDoSProtectionPlanVirtualNetworkIDs(input *[]network.SubResource) []string {
	vnetIDs := make([]string, 0)
	if input == nil {
		return vnetIDs
	}

	// if-continue is used to simplify the deeply nested if-else statement.
	for _, subRes := range *input {
		if subRes.ID != nil {
			vnetIDs = append(vnetIDs, *subRes.ID)
		}
	}

	return vnetIDs
}

func extractVnetNames(d *pluginsdk.ResourceData) (*[]string, error) {
	vnetIDs := d.Get("virtual_network_ids").([]interface{})
	vnetNames := make([]string, 0)

	for _, vnetID := range vnetIDs {
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
