package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"attached": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmPublicIPsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).publicIPClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	attachedOnly := d.Get("attached").(bool)
	resp, err := client.List(ctx, resGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response().Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure resource group %q: %v", resGroup, err)
	}

	var filteredIps []network.PublicIPAddress
	for _, element := range resp.Values() {
		nicIsAttached := element.IPConfiguration != nil

		if attachedOnly && nicIsAttached {
			filteredIps = append(filteredIps, element)
		}
	}

	var results []map[string]string
	for _, element := range filteredIps {
		flattenedIPAddress := flattenDataSourcePublicIP(element)
		results = append(results, flattenedIPAddress)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("public_ips", results)

	return nil
}

func flattenDataSourcePublicIP(input network.PublicIPAddress) map[string]string {
	output := make(map[string]string, 0)

	if input.ID != nil {
		output["id"] = *input.ID
	}

	if input.Name != nil {
		output["name"] = *input.Name
	}

	if props := input.PublicIPAddressPropertiesFormat; props != nil {
		if dns := props.DNSSettings; dns != nil {
			if fqdn := dns.Fqdn; fqdn != nil {
				output["fqdn"] = *fqdn
			}

			if label := dns.DomainNameLabel; label != nil {
				output["domain_name_label"] = *label
			}
		}

		if ip := props.IPAddress; ip != nil {
			output["ip_address"] = *ip
		}
	}

	return output
}
