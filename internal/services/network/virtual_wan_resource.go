package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualWan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualWanCreateUpdate,
		Read:   resourceVirtualWanRead,
		Update: resourceVirtualWanCreateUpdate,
		Delete: resourceVirtualWanDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualWanID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"disable_vpn_encryption": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"allow_branch_to_branch_traffic": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// TODO 3.0: remove this property
			"allow_vnet_to_vnet_traffic": {
				Type:       pluginsdk.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "this property has been removed from the API and will be removed in version 3.0 of the provider",
			},

			"office365_local_breakout_category": {
				Type:     pluginsdk.TypeString,
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Standard",
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualWanCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWanClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual WAN creation.")

	id := parse.NewVirtualWanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	location := azure.NormalizeLocation(d.Get("location").(string))
	disableVpnEncryption := d.Get("disable_vpn_encryption").(bool)
	allowBranchToBranchTraffic := d.Get("allow_branch_to_branch_traffic").(bool)
	office365LocalBreakoutCategory := d.Get("office365_local_breakout_category").(string)
	virtualWanType := d.Get("type").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, wan)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualWanRead(d, meta)
}

func resourceVirtualWanRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
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

func resourceVirtualWanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
		}
	}

	return nil
}
