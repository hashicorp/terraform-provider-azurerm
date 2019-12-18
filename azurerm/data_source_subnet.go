package azurerm

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"virtual_network_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"address_prefix": {
				Type:     schema.TypeString,
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

			"ip_configurations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
	client := meta.(*ArmClient).Network.SubnetsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	virtualNetworkName := d.Get("virtual_network_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, virtualNetworkName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Subnet %q (Virtual Network %q / Resource Group %q) was not found", name, virtualNetworkName, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Subnet %q: %+v", name, err)
	}
	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("virtual_network_name", virtualNetworkName)

	if props := resp.SubnetPropertiesFormat; props != nil {
		d.Set("address_prefix", props.AddressPrefix)

		if pe := props.PrivateEndpointNetworkPolicies; pe != nil {
			d.Set("enforce_private_link_endpoint_network_policies", strings.EqualFold("Disabled", *pe))
		}

		if ps := props.PrivateLinkServiceNetworkPolicies; ps != nil {
			d.Set("enforce_private_link_service_network_policies", strings.EqualFold("Disabled", *ps))
		}

		if props.NetworkSecurityGroup != nil {
			d.Set("network_security_group_id", props.NetworkSecurityGroup.ID)
		} else {
			d.Set("network_security_group_id", "")
		}

		if props.RouteTable != nil {
			d.Set("route_table_id", props.RouteTable.ID)
		} else {
			d.Set("route_table_id", "")
		}

		if err := d.Set("ip_configurations", flattenSubnetIPConfigurations(props.IPConfigurations)); err != nil {
			return err
		}

		if err := d.Set("service_endpoints", flattenSubnetServiceEndpoints(props.ServiceEndpoints)); err != nil {
			return err
		}
	}

	return nil
}
