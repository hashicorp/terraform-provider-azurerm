package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVnetRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"address_space": {
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

			"guid": {
				Type:     schema.TypeString,
				Computed: true,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceVnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: %s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		d.Set("guid", props.ResourceGUID)

		if as := props.AddressSpace; as != nil {
			if err := d.Set("address_space", utils.FlattenStringSlice(as.AddressPrefixes)); err != nil {
				return fmt.Errorf("error setting `address_space`: %v", err)
			}
		}

		if options := props.DhcpOptions; options != nil {
			if err := d.Set("dns_servers", utils.FlattenStringSlice(options.DNSServers)); err != nil {
				return fmt.Errorf("error setting `dns_servers`: %v", err)
			}
		}

		if err := d.Set("subnets", flattenVnetSubnetsNames(props.Subnets)); err != nil {
			return fmt.Errorf("error setting `subnets`: %v", err)
		}

		if err := d.Set("vnet_peerings", flattenVnetPeerings(props.VirtualNetworkPeerings)); err != nil {
			return fmt.Errorf("error setting `vnet_peerings`: %v", err)
		}
	}
	return nil
}

func flattenVnetSubnetsNames(input *[]network.Subnet) []interface{} {
	subnets := make([]interface{}, 0)

	if mysubnets := input; mysubnets != nil {
		for _, subnet := range *mysubnets {
			if v := subnet.Name; v != nil {
				subnets = append(subnets, *v)
			}
		}
	}
	return subnets
}

func flattenVnetPeerings(input *[]network.VirtualNetworkPeering) map[string]interface{} {
	output := make(map[string]interface{})

	if peerings := input; peerings != nil {
		for _, vnetpeering := range *peerings {
			if vnetpeering.Name == nil || vnetpeering.RemoteVirtualNetwork == nil || vnetpeering.RemoteVirtualNetwork.ID == nil {
				continue
			}

			key := *vnetpeering.Name
			value := *vnetpeering.RemoteVirtualNetwork.ID

			output[key] = value
		}
	}

	return output
}
