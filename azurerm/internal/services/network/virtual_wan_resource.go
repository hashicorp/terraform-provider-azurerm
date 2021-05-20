package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVirtualWan() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualWanCreateUpdate,
		Read:   resourceVirtualWanRead,
		Update: resourceVirtualWanCreateUpdate,
		Delete: resourceVirtualWanDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validation.StringIsNotEmpty,
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

			// TODO 3.0: remove this property
			"allow_vnet_to_vnet_traffic": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "this property has been removed from the API and will be removed in version 3.0 of the provider",
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

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Standard",
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualWanCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWanClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual WAN creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	disableVpnEncryption := d.Get("disable_vpn_encryption").(bool)
	allowBranchToBranchTraffic := d.Get("allow_branch_to_branch_traffic").(bool)
	office365LocalBreakoutCategory := d.Get("office365_local_breakout_category").(string)
	virtualWanType := d.Get("type").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
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
			Office365LocalBreakoutCategory: network.OfficeTrafficCategory(office365LocalBreakoutCategory),
			Type:                           utils.String(virtualWanType),
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

	return resourceVirtualWanRead(d, meta)
}

func resourceVirtualWanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWanClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualWanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Virtual WAN %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Virtual WAN %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualWanProperties; props != nil {
		d.Set("disable_vpn_encryption", props.DisableVpnEncryption)
		d.Set("allow_branch_to_branch_traffic", props.AllowBranchToBranchTraffic)
		d.Set("office365_local_breakout_category", props.Office365LocalBreakoutCategory)
		d.Set("allow_vnet_to_vnet_traffic", false)
		d.Set("type", props.Type)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualWanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWanClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualWanID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		// deleted outside of Terraform
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting Virtual WAN %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of Virtual WAN %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
