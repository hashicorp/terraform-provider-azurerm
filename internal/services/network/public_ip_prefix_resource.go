package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePublicIpPrefix() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePublicIpPrefixCreateUpdate,
		Read:   resourcePublicIpPrefixRead,
		Update: resourcePublicIpPrefixCreateUpdate,
		Delete: resourcePublicIpPrefixDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PublicIpPrefixID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: func() map[string]*pluginsdk.Schema {
			s := map[string]*pluginsdk.Schema{
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
					ValidateFunc: validation.IntBetween(0, 127),
				},

				"ip_version": {
					Type:             pluginsdk.TypeString,
					Optional:         true,
					Default:          string(network.IPVersionIPv4),
					ForceNew:         true,
					DiffSuppressFunc: suppress.CaseDifference,
					ValidateFunc: validation.StringInSlice([]string{
						string(network.IPVersionIPv4),
						string(network.IPVersionIPv6),
					}, true),
				},

				"ip_prefix": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"tags": tags.Schema(),
			}

			if features.ThreePointOhBeta() {
				s["zones"] = commonschema.ZonesMultipleOptionalForceNew()
			} else {
				s["availability_zone"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeString,
					Optional: true,
					// Default:  "Zone-Redundant",
					Computed: true,
					ForceNew: true,
					ConflictsWith: []string{
						"zones",
					},
					ValidateFunc: validation.StringInSlice([]string{
						"No-Zone",
						"1",
						"2",
						"3",
						"Zone-Redundant",
					}, false),
				}
				// TODO - 3.0 make Computed only
				s["zones"] = &pluginsdk.Schema{
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					ForceNew: true,
					ConflictsWith: []string{
						"availability_zone",
					},
					Deprecated: "This property has been deprecated in favour of `availability_zone` due to a breaking behavioural change in Azure: https://azure.microsoft.com/en-us/updates/zone-behavior-change/",
					MaxItems:   1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				}
			}

			return s
		}(),
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
	prefixLength := d.Get("prefix_length").(int)
	ipVersion := d.Get("ip_version").(string)
	t := d.Get("tags").(map[string]interface{})

	publicIpPrefix := network.PublicIPPrefix{
		Location: &location,
		Sku: &network.PublicIPPrefixSku{
			Name: network.PublicIPPrefixSkuName(sku),
		},
		PublicIPPrefixPropertiesFormat: &network.PublicIPPrefixPropertiesFormat{
			PrefixLength:           utils.Int32(int32(prefixLength)),
			PublicIPAddressVersion: network.IPVersion(ipVersion),
		},
		Tags: tags.Expand(t),
	}

	if features.ThreePointOhBeta() {
		zones := zones.Expand(d.Get("zones").(*schema.Set).List())
		if len(zones) > 0 {
			publicIpPrefix.Zones = &zones
		}
	} else {
		zones := &[]string{"1", "2"}
		// TODO - Remove in 3.0
		if deprecatedZonesRaw, ok := d.GetOk("zones"); ok {
			deprecatedZones := azure.ExpandZones(deprecatedZonesRaw.([]interface{}))
			if deprecatedZones != nil {
				zones = deprecatedZones
			}
		}

		if availabilityZones, ok := d.GetOk("availability_zone"); ok {
			switch availabilityZones.(string) {
			case "1", "2", "3":
				zones = &[]string{availabilityZones.(string)}
			case "Zone-Redundant":
				zones = &[]string{"1", "2"}
			case "No-Zone":
				zones = &[]string{}
			}
		}
		publicIpPrefix.Zones = zones
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

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.PublicIPPrefixeName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if features.ThreePointOhBeta() {
		d.Set("zones", zones.Flatten(resp.Zones))
	} else {
		availabilityZones := "No-Zone"
		zonesDeprecated := make([]string, 0)
		if resp.Zones != nil {
			if len(*resp.Zones) > 1 {
				availabilityZones = "Zone-Redundant"
			}
			if len(*resp.Zones) == 1 {
				zones := *resp.Zones
				availabilityZones = zones[0]
				zonesDeprecated = zones
			}
		}

		d.Set("availability_zone", availabilityZones)
		d.Set("zones", zonesDeprecated)
	}

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}
	d.Set("sku", skuName)

	if props := resp.PublicIPPrefixPropertiesFormat; props != nil {
		d.Set("prefix_length", props.PrefixLength)
		d.Set("ip_prefix", props.IPPrefix)

		if version := props.PublicIPAddressVersion; version != "" {
			d.Set("ip_version", string(props.PublicIPAddressVersion))
		}
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
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
