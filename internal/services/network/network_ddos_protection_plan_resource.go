package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DDoS protection plan creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing DDoS Protection Plan %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_ddos_protection_plan", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	vnetsToLock, err := expandNetworkDDoSProtectionPlanVnetNames(d)
	if err != nil {
		return fmt.Errorf("extracting names of Virtual Network: %+v", err)
	}

	locks.ByName(name, azureNetworkDDoSProtectionPlanResourceName)
	defer locks.UnlockByName(name, azureNetworkDDoSProtectionPlanResourceName)

	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	parameters := network.DdosProtectionPlan{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	plan, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if plan.ID == nil {
		return fmt.Errorf("Cannot read DDoS Protection Plan %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*plan.ID)

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
		vnetResourceID, err := parse.VirtualNetworkID(vnetID.(string))
		if err != nil {
			return nil, err
		}

		if !utils.SliceContainsValue(vnetNames, vnetResourceID.Name) {
			vnetNames = append(vnetNames, vnetResourceID.Name)
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
		vnetResourceID, err := parse.VirtualNetworkID(vnetID.(string))
		if err != nil {
			return nil, err
		}

		if !utils.SliceContainsValue(vnetNames, vnetResourceID.Name) {
			vnetNames = append(vnetNames, vnetResourceID.Name)
		}
	}

	return &vnetNames, nil
}
