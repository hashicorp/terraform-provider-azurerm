package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPublicIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPublicIpCreateUpdate,
		Read:   resourceArmPublicIpRead,
		Update: resourceArmPublicIpCreateUpdate,
		Delete: resourceArmPublicIpDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id, err := azure.ParseAzureResourceID(d.Id())
				if err != nil {
					return nil, err
				}
				name := id.Path["publicIPAddresses"]
				if name == "" {
					return nil, fmt.Errorf("Error parsing supplied resource id. Please check it and rerun:\n %s", d.Id())
				}
				return []*schema.ResourceData{d}, nil
			},
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"allocation_method": {
				Type: schema.TypeString,
				//Required:         true, //revert in 2.0
				Optional:      true,
				Computed:      true, // remove in 2.0
				ConflictsWith: []string{"public_ip_address_allocation"},
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Static),
					string(network.Dynamic),
				}, false),
			},

			"public_ip_address_allocation": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				StateFunc:        state.IgnoreCase,
				ConflictsWith:    []string{"allocation_method"},
				Computed:         true,
				Deprecated:       "this property has been deprecated in favor of `allocation_method` to better match the api",
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Static),
					string(network.Dynamic),
				}, true),
			},

			"ip_version": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(network.IPv4),
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.IPv4),
					string(network.IPv6),
				}, true),
			},

			"sku": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(network.PublicIPAddressSkuNameBasic),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.PublicIPAddressSkuNameBasic),
					string(network.PublicIPAddressSkuNameStandard),
				}, true),
			},

			"idle_timeout_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      4,
				ValidateFunc: validation.IntBetween(4, 30),
			},

			"domain_name_label": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.PublicIpDomainNameLabel,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"reverse_fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_ip_prefix_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPublicIpCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PublicIPsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Public IP creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	sku := d.Get("sku").(string)
	t := d.Get("tags").(map[string]interface{})
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))
	idleTimeout := d.Get("idle_timeout_in_minutes").(int)
	ipVersion := network.IPVersion(d.Get("ip_version").(string))

	ipAllocationMethod := ""
	if v, ok := d.GetOk("allocation_method"); ok {
		ipAllocationMethod = v.(string)
	} else if v, ok := d.GetOk("public_ip_address_allocation"); ok {
		ipAllocationMethod = v.(string)
	} else {
		return fmt.Errorf("Either `allocation_method` or `public_ip_address_allocation` must be specified.")
	}

	if strings.EqualFold(string(ipVersion), string(network.IPv6)) {
		if strings.EqualFold(ipAllocationMethod, "static") {
			return fmt.Errorf("Cannot specify publicIpAllocationMethod as Static for IPv6 PublicIp")
		}
	}

	if strings.EqualFold(sku, "standard") {
		if !strings.EqualFold(ipAllocationMethod, "static") {
			return fmt.Errorf("Static IP allocation must be used when creating Standard SKU public IP addresses.")
		}
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Public IP %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_public_ip", *existing.ID)
		}
	}

	publicIp := network.PublicIPAddress{
		Name:     &name,
		Location: &location,
		Sku: &network.PublicIPAddressSku{
			Name: network.PublicIPAddressSkuName(sku),
		},
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: network.IPAllocationMethod(ipAllocationMethod),
			PublicIPAddressVersion:   ipVersion,
			IdleTimeoutInMinutes:     utils.Int32(int32(idleTimeout)),
		},
		Tags:  tags.Expand(t),
		Zones: zones,
	}

	publicIpPrefixId, publicIpPrefixIdOk := d.GetOk("public_ip_prefix_id")

	if publicIpPrefixIdOk {
		publicIpPrefix := network.SubResource{}
		publicIpPrefix.ID = utils.String(publicIpPrefixId.(string))
		publicIp.PublicIPAddressPropertiesFormat.PublicIPPrefix = &publicIpPrefix
	}

	dnl, dnlOk := d.GetOk("domain_name_label")
	rfqdn, rfqdnOk := d.GetOk("reverse_fqdn")

	if dnlOk || rfqdnOk {
		dnsSettings := network.PublicIPAddressDNSSettings{}

		if rfqdnOk {
			dnsSettings.ReverseFqdn = utils.String(rfqdn.(string))
		}

		if dnlOk {
			dnsSettings.DomainNameLabel = utils.String(dnl.(string))
		}

		publicIp.PublicIPAddressPropertiesFormat.DNSSettings = &dnsSettings
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, publicIp)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Public IP %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Public IP %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Public IP %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPublicIpRead(d, meta)
}

func resourceArmPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PublicIPsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["publicIPAddresses"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Public IP %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("zones", resp.Zones)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.PublicIPAddressPropertiesFormat; props != nil {
		d.Set("public_ip_address_allocation", string(props.PublicIPAllocationMethod))
		d.Set("allocation_method", string(props.PublicIPAllocationMethod))
		d.Set("ip_version", string(props.PublicIPAddressVersion))

		if publicIpPrefix := props.PublicIPPrefix; publicIpPrefix != nil {
			d.Set("public_ip_prefix_id", publicIpPrefix.ID)
		}

		if settings := props.DNSSettings; settings != nil {
			d.Set("fqdn", settings.Fqdn)
			d.Set("reverse_fqdn", settings.ReverseFqdn)
			d.Set("domain_name_label", settings.DomainNameLabel)
		}

		d.Set("ip_address", props.IPAddress)
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPublicIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PublicIPsClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["publicIPAddresses"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Public IP %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Public IP %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
