package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	aznetapp "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetAppVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetAppVolumeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: aznetapp.ValidateNetAppPoolName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: aznetapp.ValidateNetAppAccountName,
			},

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: aznetapp.ValidateNetAppPoolName,
			},

			"creation_token": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_level": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"usage_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"export_policy_rule": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"allowed_clients": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cifs": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nfsv3": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nfsv4": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"unix_read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"unix_read_write": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.VolumeClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	poolName := d.Get("pool_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, accountName, poolName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: NetApp Volume %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading NetApp Volume %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_name", accountName)
	d.Set("pool_name", poolName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VolumeProperties; props != nil {
		d.Set("creation_token", props.CreationToken)
		d.Set("service_level", props.ServiceLevel)
		d.Set("subnet_id", props.SubnetID)

		if props.UsageThreshold != nil {
			d.Set("usage_threshold", *props.UsageThreshold/1073741824)
		}
		if props.ExportPolicy != nil {
			if err := d.Set("export_policy_rule", flattenArmNetAppVolumeExportPolicyRule(props.ExportPolicy)); err != nil {
				return fmt.Errorf("Error setting `export_policy_rule`: %+v", err)
			}
		}
	}

	return nil
}
