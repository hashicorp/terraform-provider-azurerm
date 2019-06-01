package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const azureNetworkDDoSProtectionPlanResourceName = "azurerm_network_ddos_protection_plan"

func resourceArmNetworkDDoSProtectionPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkDDoSProtectionPlanCreateUpdate,
		Read:   resourceArmNetworkDDoSProtectionPlanRead,
		Update: resourceArmNetworkDDoSProtectionPlanCreateUpdate,
		Delete: resourceArmNetworkDDoSProtectionPlanDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"virtual_network_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkDDoSProtectionPlanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ddosProtectionPlanClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DDoS protection plan creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing DDoS Protection Plan %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_ddos_protection_plan", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	vnetsToLock, err := expandArmNetworkDDoSProtectionPlanVnetNames(d)
	if err != nil {
		return fmt.Errorf("Error extracting names of Virtual Network: %+v", err)
	}

	azureRMLockByName(name, azureNetworkDDoSProtectionPlanResourceName)
	defer azureRMUnlockByName(name, azureNetworkDDoSProtectionPlanResourceName)

	azureRMLockMultipleByName(vnetsToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(vnetsToLock, virtualNetworkResourceName)

	parameters := network.DdosProtectionPlan{
		Location: &location,
		Tags:     expandTags(tags),
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

	return resourceArmNetworkDDoSProtectionPlanRead(d, meta)
}

func resourceArmNetworkDDoSProtectionPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ddosProtectionPlanClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["NetworkDDoSProtectionPlans"]

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
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := plan.DdosProtectionPlanPropertiesFormat; props != nil {
		vNetIDs := flattenArmNetworkDDoSProtectionPlanVirtualNetworkIDs(props.VirtualNetworks)
		if err := d.Set("virtual_network_ids", vNetIDs); err != nil {
			return fmt.Errorf("Error setting `virtual_network_ids`: %+v", err)
		}
	}

	flattenAndSetTags(d, plan.Tags)

	return nil
}

func resourceArmNetworkDDoSProtectionPlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ddosProtectionPlanClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["NetworkDDoSProtectionPlans"]

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

	azureRMLockByName(name, azureNetworkDDoSProtectionPlanResourceName)
	defer azureRMUnlockByName(name, azureNetworkDDoSProtectionPlanResourceName)

	azureRMLockMultipleByName(vnetsToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(vnetsToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandArmNetworkDDoSProtectionPlanVnetNames(d *schema.ResourceData) (*[]string, error) {
	vnetIDs := d.Get("virtual_network_ids").([]interface{})
	vnetNames := make([]string, 0)

	for _, vnetID := range vnetIDs {
		vnetResourceID, err := parseAzureResourceID(vnetID.(string))
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

func flattenArmNetworkDDoSProtectionPlanVirtualNetworkIDs(input *[]network.SubResource) []string {
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
