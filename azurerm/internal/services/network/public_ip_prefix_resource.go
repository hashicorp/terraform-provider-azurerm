package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePublicIpPrefix() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePublicIpPrefixCreateUpdate,
		Read:   resourcePublicIpPrefixRead,
		Update: resourcePublicIpPrefixCreateUpdate,
		Delete: resourcePublicIpPrefixDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.PublicIPPrefixSkuNameStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.PublicIPPrefixSkuNameStandard),
				}, false),
			},

			"prefix_length": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      28,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 31),
			},

			"ip_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func resourcePublicIpPrefixCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	// TODO: remove in 3.0
	// to address this breaking change : https://azure.microsoft.com/en-us/updates/zone-behavior-change/, by setting a location's available zone list as default value to keep the behavior unchanged,
	// to create a non-zonal resource, user can set zones to ["no_zone"]
	if zones == nil || len(*zones) == 0 {
		allZones, err := getZones(ctx, meta.(*clients.Client).Resource.ResourceProvidersClient, "publicIPPrefixes", location)
		if err != nil {
			return fmt.Errorf("failed to get available zones for resourceType: publicIPAddresses, location: %s:%+v", location, err)
		} else {
			zones = allZones
		}
	} else if (*zones)[0] == "no_zone" {
		zones = nil
	}

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

func resourcePublicIpPrefixRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	// TODO: remove in 3.0
	zones := make([]string, 0)
	if resp.Location != nil && resp.Sku != nil && strings.EqualFold(string(resp.Sku.Name), "standard") {
		allZones, err := getZones(ctx, meta.(*clients.Client).Resource.ResourceProvidersClient, "publicIPPrefixes", *resp.Location)
		if err != nil {
			return fmt.Errorf("failed to get available zones for resourceType: publicIPPrefixes, location: %s:%+v", *resp.Location, err)
		}
		switch {
		case allZones == nil || len(*allZones) == 0:
			zones = make([]string, 0)
		case resp.Zones == nil || len(*resp.Zones) == 0:
			zones = append(zones, "no_zone")
		case len(*resp.Zones) == 1:
			zones = append(zones, (*resp.Zones)[0])
		default:
			zones = make([]string, 0)
		}
	}
	d.Set("zones", &zones)
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

func resourcePublicIpPrefixDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
