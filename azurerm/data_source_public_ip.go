package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"allocation_method": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_name_label": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"reverse_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaZonesComputed(),

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PublicIPsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Public IP %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("Error making Read request on Azure public ip %s: %s", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("zones", resp.Zones)

	//ensure values are at least set to "", d.Set() is a noop on a nil
	//there must be a better way...
	d.Set("location", "")
	d.Set("sku", "")
	d.Set("fqdn", "")
	d.Set("reverse_fqdn", "")
	d.Set("domain_name_label", "")
	d.Set("allocation_method", "")
	d.Set("ip_address", "")
	d.Set("ip_version", "")
	d.Set("idle_timeout_in_minutes", 0)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.PublicIPAddressPropertiesFormat; props != nil {
		if dnsSettings := props.DNSSettings; dnsSettings != nil {
			d.Set("fqdn", dnsSettings.Fqdn)
			d.Set("reverse_fqdn", dnsSettings.ReverseFqdn)
			d.Set("domain_name_label", dnsSettings.DomainNameLabel)
		}

		d.Set("allocation_method", string(props.PublicIPAllocationMethod))
		d.Set("ip_address", props.IPAddress)
		d.Set("ip_version", string(props.PublicIPAddressVersion))
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
