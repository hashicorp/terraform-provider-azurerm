package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSubnetRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"virtual_network_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"address_prefixes": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"network_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"enforce_private_link_endpoint_network_policies": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enforce_private_link_service_network_policies": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSubnetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: %s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("virtual_network_name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.SubnetPropertiesFormat; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		if props.AddressPrefixes == nil {
			if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
				d.Set("address_prefixes", []string{*props.AddressPrefix})
			} else {
				d.Set("address_prefixes", []string{})
			}
		} else {
			d.Set("address_prefixes", utils.FlattenStringSlice(props.AddressPrefixes))
		}

		d.Set("enforce_private_link_endpoint_network_policies", flattenSubnetPrivateLinkNetworkPolicy(props.PrivateEndpointNetworkPolicies))
		d.Set("enforce_private_link_service_network_policies", flattenSubnetPrivateLinkNetworkPolicy(props.PrivateLinkServiceNetworkPolicies))

		networkSecurityGroupId := ""
		if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.ID != nil {
			networkSecurityGroupId = *props.NetworkSecurityGroup.ID
		}
		d.Set("network_security_group_id", networkSecurityGroupId)

		routeTableId := ""
		if props.RouteTable != nil && props.RouteTable.ID != nil {
			routeTableId = *props.RouteTable.ID
		}
		d.Set("route_table_id", routeTableId)

		if err := d.Set("service_endpoints", flattenSubnetServiceEndpoints(props.ServiceEndpoints)); err != nil {
			return fmt.Errorf("setting `service_endpoints`: %+v", err)
		}
	}

	return nil
}
