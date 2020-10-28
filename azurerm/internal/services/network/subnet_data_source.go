package network

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubnetRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// TODO 3.0 - remove this property
			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),

			// TODO 3.0 - remove this property
			"virtual_network_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "This property is deprecated in favor of `virtual_network_id`",
				ExactlyOneOf: []string{"virtual_network_name", "virtual_network_id"},
				RequiredWith: []string{"resource_group_name"},
			},

			"virtual_network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"virtual_network_name", "virtual_network_id"},
				ConflictsWith: []string{"resource_group_name"},
			},

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

func dataSourceArmSubnetRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	var vnetId parse.VirtualNetworkId
	if d.Get("virtual_network_id") != "" {
		vnetIdPtr, err := parse.VirtualNetworkID(d.Get("virtual_network_id").(string))
		if err != nil {
			return err
		}
		vnetId = *vnetIdPtr
	} else {
		vnetId = parse.NewVirtualNetworkID(d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string))
	}

	resp, err := client.Get(ctx, vnetId.ResourceGroup, vnetId.Name, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Subnet %q (Virtual Network %q / Resource Group %q) was not found", name, vnetId.Name, vnetId.ResourceGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Subnet %q: %+v", name, err)
	}

	id, err := parse.SubnetID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	d.Set("name", name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("virtual_network_name", id.VirtualNetworkName)
	d.Set("virtual_network_id", parse.NewVirtualNetworkID(id.ResourceGroup, id.VirtualNetworkName).ID(subscriptionId))

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
			return fmt.Errorf("Error setting `service_endpoints`: %+v", err)
		}
	}

	return nil
}
