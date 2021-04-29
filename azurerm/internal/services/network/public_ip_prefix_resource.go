package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
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

func resourcePublicIpPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIpPrefixCreateUpdate,
		Read:   resourcePublicIpPrefixRead,
		Update: resourcePublicIpPrefixCreateUpdate,
		Delete: resourcePublicIpPrefixDelete,
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.Standard),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Standard),
				}, false),
			},

			"prefix_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      28,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 31),
			},

			"ip_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func resourcePublicIpPrefixCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Public IP Prefix creation.")

	id := parse.NewPublicIpPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.PublicIPPrefixeName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_public_ip_prefix", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	prefix_length := d.Get("prefix_length").(int)
	t := d.Get("tags").(map[string]interface{})
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	publicIpPrefix := network.PublicIPPrefix{
		Location: &location,
		Sku: &network.PublicIPPrefixSku{
			Name: network.PublicIPPrefixSkuName(sku),
		},
		PublicIPPrefixPropertiesFormat: &network.PublicIPPrefixPropertiesFormat{
			PrefixLength: utils.Int32(int32(prefix_length)),
		},
		Tags:  tags.Expand(t),
		Zones: zones,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.PublicIPPrefixeName, publicIpPrefix)
	if err != nil {
		return fmt.Errorf("creating/Updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePublicIpPrefixRead(d, meta)
}

func resourcePublicIpPrefixRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PublicIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.PublicIPPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.Set("name", id.PublicIPPrefixeName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("zones", resp.Zones)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.PublicIPPrefixPropertiesFormat; props != nil {
		d.Set("prefix_length", props.PrefixLength)
		d.Set("ip_prefix", props.IPPrefix)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePublicIpPrefixDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PublicIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.PublicIPPrefixeName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
