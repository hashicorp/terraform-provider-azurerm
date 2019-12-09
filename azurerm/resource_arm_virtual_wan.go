package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	networkSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualWan() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualWanCreateUpdate,
		Read:   resourceArmVirtualWanRead,
		Update: resourceArmVirtualWanCreateUpdate,
		Delete: resourceArmVirtualWanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"disable_vpn_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allow_branch_to_branch_traffic": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allow_vnet_to_vnet_traffic": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"office365_local_breakout_category": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.OfficeTrafficCategoryAll),
					string(network.OfficeTrafficCategoryNone),
					string(network.OfficeTrafficCategoryOptimize),
					string(network.OfficeTrafficCategoryOptimizeAndAllow),
				}, false),
				Default: string(network.OfficeTrafficCategoryNone),
			},

			"tags": tags.Schema(),

			// Remove in 2.0
			"security_provider_name": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This field has been removed by Azure and will be removed in version 2.0 of the Azure Provider",
			},
		},
	}
}

func resourceArmVirtualWanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualWanClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual WAN creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	disableVpnEncryption := d.Get("disable_vpn_encryption").(bool)
	allowBranchToBranchTraffic := d.Get("allow_branch_to_branch_traffic").(bool)
	allowVnetToVnetTraffic := d.Get("allow_vnet_to_vnet_traffic").(bool)
	office365LocalBreakoutCategory := d.Get("office365_local_breakout_category").(string)
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_wan", *existing.ID)
		}
	}

	wan := network.VirtualWAN{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		VirtualWanProperties: &network.VirtualWanProperties{
			DisableVpnEncryption:           utils.Bool(disableVpnEncryption),
			AllowBranchToBranchTraffic:     utils.Bool(allowBranchToBranchTraffic),
			AllowVnetToVnetTraffic:         utils.Bool(allowVnetToVnetTraffic),
			Office365LocalBreakoutCategory: network.OfficeTrafficCategory(office365LocalBreakoutCategory),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, wan)
	if err != nil {
		return fmt.Errorf("Error creating/updating Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual WAN %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualWanRead(d, meta)
}

func resourceArmVirtualWanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualWanClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := networkSvc.ParseVirtualWanID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.Base.ResourceGroup
	name := id.Name

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual WAN %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualWanProperties; props != nil {
		d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
		d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
		d.Set("allow_vnet_to_vnet_traffic", props.AllowVnetToVnetTraffic)
		d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualWanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.VirtualWanClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := networkSvc.ParseVirtualWanID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.Base.ResourceGroup
	name := id.Name

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		// deleted outside of Terraform
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of Virtual WAN %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
