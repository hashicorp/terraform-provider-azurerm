package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/schema"
	//"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVnetRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"address_spaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmVnetRead(d *schema.ResourceData, meta interface{}) error {
	vnetClient := meta.(*ArmClient).vnetClient

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := vnetClient.Get(resGroup, name, "")
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
		}
		return fmt.Errorf("Error making Read request on Azure public ip %s: %s", name, err)
	}

	d.SetId(*resp.ID)

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		address_spaces := flattenVnetAddressPrefixes(props.AddressSpace.AddressPrefixes)
		if err := d.Set("address_spaces", address_spaces); err != nil {
			return err
		}

		dns_servers := flattenVnetAddressPrefixes(props.DhcpOptions.DNSServers)
		if err := d.Set("dns_servers", dns_servers); err != nil {
			return err
		}

		subnets := flattenVnetSubnetsAddressSpace(props.Subnets)
		if err := d.Set("subnets", subnets); err != nil {
			return err
		}
	}

	return nil
}

func flattenVnetAddressPrefixes(input *[]string) []interface{} {
	prefixes := make([]interface{}, 0)

	for _, prefix := range *input {
		prefixes = append(prefixes, prefix)
	}

	return prefixes
}

func flattenVnetSubnetsAddressSpace(input *[]network.Subnet) []interface{} {
	subnets := make([]interface{}, 0)

	for _, subnet := range *input {
		subnets = append(subnets, *subnet.Name)
	}

	return subnets
}
