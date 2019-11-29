package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const azureDDoSProtectionPlanResourceName = "azurerm_ddos_protection_plan"

func resourceArmDDoSProtectionPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDDoSProtectionPlanCreateUpdate,
		Read:   resourceArmDDoSProtectionPlanRead,
		Update: resourceArmDDoSProtectionPlanCreateUpdate,
		Delete: resourceArmDDoSProtectionPlanDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		DeprecationMessage: `The 'azurerm_ddos_protection_plan' resource is deprecated in favour of the renamed version 'azurerm_network_ddos_protection_plan'.

Information on migrating to the renamed resource can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html

As such the existing 'azurerm_ddos_protection_plan' resource is deprecated and will be removed in the next major version of the AzureRM Provider (2.0).
`,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_network_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDDoSProtectionPlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.DDOSProtectionPlansClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DDoS protection plan creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DDoS Protection Plan %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_ddos_protection_plan", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	vnetsToLock, err := extractVnetNames(d)
	if err != nil {
		return fmt.Errorf("Error extracting names of Virtual Network: %+v", err)
	}

	locks.ByName(name, azureDDoSProtectionPlanResourceName)
	defer locks.UnlockByName(name, azureDDoSProtectionPlanResourceName)

	locks.MultipleByName(vnetsToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, virtualNetworkResourceName)

	parameters := network.DdosProtectionPlan{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	plan, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if plan.ID == nil {
		return fmt.Errorf("Cannot read DDoS Protection Plan %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*plan.ID)

	return resourceArmDDoSProtectionPlanRead(d, meta)
}

func resourceArmDDoSProtectionPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.DDOSProtectionPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ddosProtectionPlans"]

	plan, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(plan.Response) {
			log.Printf("[DEBUG] DDoS Protection Plan %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", plan.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := plan.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := plan.DdosProtectionPlanPropertiesFormat; props != nil {
		vNetIDs := flattenArmVirtualNetworkIDs(props.VirtualNetworks)
		if err := d.Set("virtual_network_ids", vNetIDs); err != nil {
			return fmt.Errorf("Error setting `virtual_network_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, plan.Tags)
}

func resourceArmDDoSProtectionPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.DDOSProtectionPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ddosProtectionPlans"]

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] DDoS Protection Plan %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	vnetsToLock, err := extractVnetNames(d)
	if err != nil {
		return fmt.Errorf("Error extracting names of Virtual Network: %+v", err)
	}

	locks.ByName(name, azureDDoSProtectionPlanResourceName)
	defer locks.UnlockByName(name, azureDDoSProtectionPlanResourceName)

	locks.MultipleByName(vnetsToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func extractVnetNames(d *schema.ResourceData) (*[]string, error) {
	vnetIDs := d.Get("virtual_network_ids").([]interface{})
	vnetNames := make([]string, 0)

	for _, vnetID := range vnetIDs {
		vnetResourceID, err := azure.ParseAzureResourceID(vnetID.(string))
		if err != nil {
			return nil, err
		}

		vnetName := vnetResourceID.Path["virtualNetworks"]

		if !sliceContainsValue(vnetNames, vnetName) {
			vnetNames = append(vnetNames, vnetName)
		}
	}

	return &vnetNames, nil
}

func flattenArmVirtualNetworkIDs(input *[]network.SubResource) []string {
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
