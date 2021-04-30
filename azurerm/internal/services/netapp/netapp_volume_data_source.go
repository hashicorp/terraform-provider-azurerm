package netapp

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/netapp/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceNetAppVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetAppVolumeRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},

			"mount_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"volume_path": {
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

			"storage_quota_in_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"data_protection_replication": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"remote_volume_location": azure.SchemaLocationForDataSource(),

						"remote_volume_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"replication_schedule": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "This property is not in use and will be removed in version 3.0 of the provider. Please use `replication_frequency` instead",
						},

						"replication_frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetAppVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Error retrieving NetApp Volume %q (Resource Group %q): ID was nil or empty", name, resourceGroup)
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
		d.Set("volume_path", props.CreationToken)
		d.Set("service_level", props.ServiceLevel)
		d.Set("subnet_id", props.SubnetID)

		protocolTypes := make([]string, 0)
		if prtclTypes := props.ProtocolTypes; prtclTypes != nil {
			protocolTypes = append(protocolTypes, *prtclTypes...)
		}
		d.Set("protocols", protocolTypes)

		if props.UsageThreshold != nil {
			d.Set("storage_quota_in_gb", *props.UsageThreshold/1073741824)
		}
		if err := d.Set("mount_ip_addresses", flattenNetAppVolumeMountIPAddresses(props.MountTargets)); err != nil {
			return fmt.Errorf("setting `mount_ip_addresses`: %+v", err)
		}
		if err := d.Set("data_protection_replication", flattenNetAppVolumeDataProtectionReplication(props.DataProtection)); err != nil {
			return fmt.Errorf("setting `data_protection_replication`: %+v", err)
		}
	}

	return nil
}
