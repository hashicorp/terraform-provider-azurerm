package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

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

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func dataSourceArmPublicIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).publicIPClient
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

	//ensure values are at least set to "", d.Set() is a noop on a nil
	d.Set("fqdn", "")
	d.Set("domain_name_label", "")
	d.Set("ip_address", "")
	d.Set("idle_timeout_in_minutes", 0)

	if props := resp.PublicIPAddressPropertiesFormat; props != nil {
		if dnsSettings := props.DNSSettings; dnsSettings != nil {
			d.Set("fqdn", dnsSettings.Fqdn)
			d.Set("domain_name_label", dnsSettings.DomainNameLabel)
		}

		d.Set("ip_address", props.IPAddress)
		d.Set("ip_version", string(props.PublicIPAddressVersion))
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}
