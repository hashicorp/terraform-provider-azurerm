// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing MySQL Flexible Server.",
		Read:        dataSourceMysqlFlexibleServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the MySQL Flexible Server.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FlexibleServerName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"administrator_login": {
				Description: "The Administrator login of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"backup_retention_days": {
				Description: "The backup retention days of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"delegated_subnet_id": {
				Description: "The ID of the Virtual Network subnet the MySQL Flexible Server is created in.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"geo_redundant_backup_enabled": {
				Description: "Whether geo redundant backup is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"high_availability": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Description: "The high availability mode of the MySQL Flexible Server.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"standby_availability_zone": {
							Description: "The Availability Zone of the standby Flexible Server.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"maintenance_window": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Description: "The day of week of the maintenance window.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"start_hour": {
							Description: "The start hour of the maintenance window.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"start_minute": {
							Description: "The start minute of the maintenance window.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
					},
				},
			},

			"restore_point_in_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_dns_zone_id": {
				Description: "The ID of the Private DNS zone of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"replication_role": {
				Description: "The replication role of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"sku_name": {
				Description: "The SKU Name of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"auto_grow_enabled": {
							Description: "Whether Storage Auto Grow is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"iops": {
							Description: "The storage IOPS of the MySQL Flexible Server.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"size_gb": {
							Description: "The max storage allowed for the MySQL Flexible Server.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},
						"io_scaling_enabled": {
							Description: "Whether IOPS is scaled automatically.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"version": {
				Description: "The version of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"zone": commonschema.ZoneSingleComputed(),

			"fqdn": {
				Description: "The fully qualified domain name of the MySQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"public_network_access_enabled": {
				Description: "Whether the public network access is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"replica_capacity": {
				Description: "The maximum number of replicas that a primary MySQL Flexible Server can have.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMysqlFlexibleServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := servers.NewFlexibleServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("zone", props.AvailabilityZone)
			d.Set("version", pointer.FromEnum(props.Version))
			d.Set("fqdn", props.FullyQualifiedDomainName)

			if network := props.Network; network != nil {
				d.Set("public_network_access_enabled", *network.PublicNetworkAccess == servers.EnableStatusEnumEnabled)
				d.Set("delegated_subnet_id", network.DelegatedSubnetResourceId)
				d.Set("private_dns_zone_id", network.PrivateDnsZoneResourceId)
			}

			if err := d.Set("maintenance_window", flattenDataSourceArmServerMaintenanceWindow(props.MaintenanceWindow)); err != nil {
				return fmt.Errorf("setting `maintenance_window`: %+v", err)
			}

			if err := d.Set("storage", flattenDataSourceArmServerStorage(props.Storage)); err != nil {
				return fmt.Errorf("setting `storage`: %+v", err)
			}

			if backup := props.Backup; backup != nil {
				d.Set("backup_retention_days", backup.BackupRetentionDays)
				d.Set("geo_redundant_backup_enabled", *backup.GeoRedundantBackup == servers.EnableStatusEnumEnabled)
			}

			if err := d.Set("high_availability", flattenDataSourceFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
				return fmt.Errorf("setting `high_availability`: %+v", err)
			}
			d.Set("replication_role", pointer.FromEnum(props.ReplicationRole))
			d.Set("replica_capacity", props.ReplicaCapacity)
		}

		sku, err := flattenDataSourceFlexibleServerSku(model.Sku)
		if err != nil {
			return fmt.Errorf("flattening `sku_name` for %q: %v", id, err)
		}

		d.Set("sku_name", sku)

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func flattenDataSourceArmServerStorage(storage *servers.Storage) []any {
	if storage == nil {
		return []any{}
	}

	var size, iops int64
	if storage.StorageSizeGB != nil {
		size = *storage.StorageSizeGB
	}

	if storage.Iops != nil {
		iops = *storage.Iops
	}

	return []any{
		map[string]any{
			"size_gb":            size,
			"iops":               iops,
			"auto_grow_enabled":  *storage.AutoGrow == servers.EnableStatusEnumEnabled,
			"io_scaling_enabled": *storage.AutoIoScaling == servers.EnableStatusEnumEnabled,
		},
	}
}

func flattenDataSourceFlexibleServerSku(sku *servers.MySQLServerSku) (string, error) {
	if sku == nil || sku.Name == "" || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case servers.ServerSkuTierBurstable:
		tier = "B"
	case servers.ServerSkuTierGeneralPurpose:
		tier = "GP"
	case servers.ServerSkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, sku.Name}, "_"), nil
}

func flattenDataSourceArmServerMaintenanceWindow(input *servers.MaintenanceWindow) []any {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == string(ServerMaintenanceWindowDisabled) {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"day_of_week":  pointer.From(input.DayOfWeek),
			"start_hour":   pointer.From(input.StartHour),
			"start_minute": pointer.From(input.StartMinute),
		},
	}
}

func flattenDataSourceFlexibleServerHighAvailability(ha *servers.HighAvailability) []any {
	if ha == nil || *ha.Mode == servers.HighAvailabilityModeDisabled {
		return []any{}
	}

	return []any{
		map[string]any{
			"mode":                      string(*ha.Mode),
			"standby_availability_zone": pointer.From(ha.StandbyAvailabilityZone),
		},
	}
}
