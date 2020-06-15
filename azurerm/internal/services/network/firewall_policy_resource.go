package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var azureFirewallPolicyResourceName = "azurerm_firewall_policy"

func resourceArmFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFirewallPolicyCreateUpdate,
		Read:   resourceArmFirewallPolicyRead,
		Update: resourceArmFirewallPolicyCreateUpdate,
		Delete: resourceArmFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"base_policy_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"threat_intel_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallThreatIntelModeAlert),
					string(network.AzureFirewallThreatIntelModeDeny),
					string(network.AzureFirewallThreatIntelModeOff),
				}, false),
			},

			"intrusion_system_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.FirewallPolicyIntrusionSystemModeEnabled),
					string(network.FirewallPolicyIntrusionSystemModeDisabled),
				}, false),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmFirewallPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall Policy creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Firewall Policy %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_firewall_policy", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	basePolicyId := d.Get("base_policy_id").(string)
	threatIntelMode := network.AzureFirewallThreatIntelMode(d.Get("threat_intel_mode").(string))
	intrusionSystemMode := network.FirewallPolicyIntrusionSystemMode(d.Get("intrusion_system_mode").(string))

	locks.ByName(name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(name, azureFirewallPolicyResourceName)

	parameters := network.FirewallPolicy{
		Location: &location,
		Tags:     tags.Expand(t),
		FirewallPolicyPropertiesFormat: &network.FirewallPolicyPropertiesFormat{
			ThreatIntelMode:     threatIntelMode,
			IntrusionSystemMode: intrusionSystemMode,
		},
	}

	if basePolicyId != "" {
		parameters.BasePolicy = &network.SubResource{
			ID: &basePolicyId,
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Azure Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Azure Firewall Policy %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmFirewallPolicyRead(d, meta)
}

func resourceArmFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["firewallPolicies"]

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Firewall Policy %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if basePolicy := read.BasePolicy; basePolicy != nil {
		d.Set("base_policy_id", basePolicy.ID)
	}
	d.Set("threat_intel_mode", read.ThreatIntelMode)
	d.Set("intrusion_system_mode", read.IntrusionSystemMode)

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceArmFirewallPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FirewallPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["firewallPolicies"]

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall Policy %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	locks.ByName(name, azureFirewallPolicyResourceName)
	defer locks.UnlockByName(name, azureFirewallPolicyResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Azure Firewall Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}
