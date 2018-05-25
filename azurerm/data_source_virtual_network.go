package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVnetRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

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

			"vnet_peerings": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceArmVnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Network %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on Virtual Network %q (resource group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		addressSpaces := flattenVnetAddressPrefixes(props.AddressSpace.AddressPrefixes)
		if err := d.Set("address_spaces", addressSpaces); err != nil {
			return err
		}

		if options := props.DhcpOptions; options != nil {
			dnsServers := flattenVnetAddressPrefixes(options.DNSServers)
			if err := d.Set("dns_servers", dnsServers); err != nil {
				return err
			}
		}

		subnets := flattenVnetSubnetsNames(props.Subnets)
		if err := d.Set("subnets", subnets); err != nil {
			return err
		}

		vnetPeerings := flattenVnetPeerings(props.VirtualNetworkPeerings)
		if err := d.Set("vnet_peerings", vnetPeerings); err != nil {
			return err
		}
	}
	return nil
}

func flattenVnetAddressPrefixes(input *[]string) []interface{} {
	prefixes := make([]interface{}, 0)

	if myprefixes := input; myprefixes != nil {
		for _, prefix := range *myprefixes {
			prefixes = append(prefixes, prefix)
		}
	}
	return prefixes
}

func flattenVnetSubnetsNames(input *[]network.Subnet) []interface{} {
	subnets := make([]interface{}, 0)

	if mysubnets := input; mysubnets != nil {
		for _, subnet := range *mysubnets {
			subnets = append(subnets, *subnet.Name)
		}
	}
	return subnets
}

func flattenVnetPeerings(input *[]network.VirtualNetworkPeering) map[string]interface{} {
	output := make(map[string]interface{}, 0)

	if peerings := input; peerings != nil {
		for _, vnetpeering := range *peerings {
			key := *vnetpeering.Name
			value := *vnetpeering.RemoteVirtualNetwork.ID

			output[key] = value

		}
	}
	return output
}
