package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDevTestVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDevTestVnetRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"lab_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"unique_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"allowed_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lab_subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"subnet_overrides": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lab_subnet_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_in_vm_creation_permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_public_ip_address_permission": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_network_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmDevTestVnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	labName := d.Get("lab_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, labName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Network %q in Dev Test Lab %q (Resource Group %q) was not found", name, labName, resGroup)
		}

		return fmt.Errorf("Error making Read request on Virtual Network %q in Dev Test Lab %q (Resource Group %q): %+v", name, labName, resGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Virtual Network %q in Dev Test Lab %q (Resource Group %q): %+v", name, labName, resGroup, err)
	}
	d.SetId(*resp.ID)

	if props := resp.VirtualNetworkProperties; props != nil {
		if as := props.AllowedSubnets; as != nil {
			if err := d.Set("allowed_subnets", flattenDevTestVirtualNetworkAllowedSubnets(as)); err != nil {
				return fmt.Errorf("error setting `allowed_subnets`: %v", err)
			}
		}
		if so := props.SubnetOverrides; so != nil {
			if err := d.Set("subnet_overrides", flattenDevTestVirtualNetworkSubnetOverrides(so)); err != nil {
				return fmt.Errorf("error setting `subnet_overrides`: %v", err)
			}
		}
		d.Set("unique_identifier", props.UniqueIdentifier)
	}
	return nil
}

func flattenDevTestVirtualNetworkAllowedSubnets(input *[]dtl.Subnet) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	for _, v := range *input {
		allowedSubnet := make(map[string]interface{})

		allowedSubnet["allow_public_ip"] = string(v.AllowPublicIP)

		if resourceID := v.ResourceID; resourceID != nil {
			allowedSubnet["resource_id"] = *resourceID
		}

		if labSubnetName := v.LabSubnetName; labSubnetName != nil {
			allowedSubnet["lab_subnet_name"] = *labSubnetName
		}

		result = append(result, allowedSubnet)
	}

	return result
}

func flattenDevTestVirtualNetworkSubnetOverrides(input *[]dtl.SubnetOverride) []interface{} {
	result := make([]interface{}, 0)

	if input == nil {
		return result
	}

	for _, v := range *input {
		subnetOverride := make(map[string]interface{})
		if v.LabSubnetName != nil {
			subnetOverride["lab_subnet_name"] = *v.LabSubnetName
		}
		if v.ResourceID != nil {
			subnetOverride["resource_id"] = *v.ResourceID
		}

		subnetOverride["use_public_ip_address_permission"] = string(v.UsePublicIPAddressPermission)
		subnetOverride["use_in_vm_creation_permission"] = string(v.UseInVMCreationPermission)

		if v.VirtualNetworkPoolName != nil {
			subnetOverride["virtual_network_pool_name"] = *v.VirtualNetworkPoolName
		}

		result = append(result, subnetOverride)
	}

	return result
}
