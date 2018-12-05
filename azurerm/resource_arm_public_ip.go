package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPublicIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPublicIpCreate,
		Read:   resourceArmPublicIpRead,
		Update: resourceArmPublicIpCreate,
		Delete: resourceArmPublicIpDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				id, err := parseAzureResourceID(d.Id())
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"zones": singleZonesSchema(),

			//should this perhaps be allocation_method? (yes i think so)
			"public_ip_address_allocation": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				StateFunc:        ignoreCaseStateFunc,
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

			"reverse_fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPublicIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).publicIPClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Public IP creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location"))
	resGroup := d.Get("resource_group_name").(string)
	sku := network.PublicIPAddressSku{
		Name: network.PublicIPAddressSkuName(d.Get("sku").(string)),
	}
	tags := d.Get("tags").(map[string]interface{})
	zones := expandZones(d.Get("zones").([]interface{}))

	idleTimeout := d.Get("idle_timeout_in_minutes").(int)
	ipAllocationMethod := network.IPAllocationMethod(d.Get("public_ip_address_allocation").(string))
	ipVersion := network.IPVersion(d.Get("ip_version").(string))

	if strings.EqualFold(string(ipVersion), string(network.IPv6)) {
		if strings.EqualFold(string(ipAllocationMethod), "static") {
			return fmt.Errorf("Cannot specify publicIpAllocationMethod as Static for IPv6 PublicIp")
		}
	}

	if strings.ToLower(string(sku.Name)) == "standard" {
		if strings.ToLower(string(ipAllocationMethod)) != "static" {
			return fmt.Errorf("Static IP allocation must be used when creating Standard SKU public IP addresses.")
		}
	}

	publicIp := network.PublicIPAddress{
		Name:     &name,
		Location: &location,
		Sku:      &sku,
		PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: ipAllocationMethod,
			PublicIPAddressVersion:   ipVersion,
			IdleTimeoutInMinutes:     utils.Int32(int32(idleTimeout)),
		},
		Tags:  expandTags(tags),
		Zones: zones,
	}

	dnl, dnlOk := d.GetOk("domain_name_label")
	rfqdn, rfqdnOk := d.GetOk("reverse_fqdn")

	if dnlOk || rfqdnOk {
		dnsSettings := network.PublicIPAddressDNSSettings{}

		if rfqdnOk {
			reverseFqdn := rfqdn.(string)
			dnsSettings.ReverseFqdn = &reverseFqdn
		}

		if dnlOk {
			domainNameLabel := dnl.(string)
			dnsSettings.DomainNameLabel = &domainNameLabel
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
	client := meta.(*ArmClient).publicIPClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
	if err := azure.FlattenAndSetLocation(d, resp.Location); err != nil {
		return err
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.PublicIPAddressPropertiesFormat; props != nil {
		d.Set("public_ip_address_allocation", string(props.PublicIPAllocationMethod))
		d.Set("ip_version", string(props.PublicIPAddressVersion))

		if settings := props.DNSSettings; settings != nil {
			d.Set("fqdn", settings.Fqdn)
			d.Set("domain_name_label", settings.DomainNameLabel)
		}

		d.Set("ip_address", props.IPAddress)
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmPublicIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).publicIPClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
