// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMysqlFlexibleServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMysqlFlexibleServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FlexibleServerName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"backup_retention_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"delegated_subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"high_availability": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"standby_availability_zone": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"start_hour": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"start_minute": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"restore_point_in_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_dns_zone_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"replication_role": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"auto_grow_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"iops": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"size_gb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"io_scaling_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zone": commonschema.ZoneSingleComputed(),

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"replica_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMysqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("version", string(pointer.From(props.Version)))
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
			d.Set("replication_role", string(pointer.From(props.ReplicationRole)))
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

func flattenDataSourceArmServerStorage(storage *servers.Storage) []interface{} {
	if storage == nil {
		return []interface{}{}
	}

	var size, iops int64
	if storage.StorageSizeGB != nil {
		size = *storage.StorageSizeGB
	}

	if storage.Iops != nil {
		iops = *storage.Iops
	}

	return []interface{}{
		map[string]interface{}{
			"size_gb":            size,
			"iops":               iops,
			"auto_grow_enabled":  *storage.AutoGrow == servers.EnableStatusEnumEnabled,
			"io_scaling_enabled": *storage.AutoIoScaling == servers.EnableStatusEnumEnabled,
		},
	}
}

func flattenDataSourceFlexibleServerSku(sku *servers.Sku) (string, error) {
	if sku == nil || sku.Name == "" || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case servers.SkuTierBurstable:
		tier = "B"
	case servers.SkuTierGeneralPurpose:
		tier = "GP"
	case servers.SkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, sku.Name}, "_"), nil
}

func flattenDataSourceArmServerMaintenanceWindow(input *servers.MaintenanceWindow) []interface{} {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == string(ServerMaintenanceWindowDisabled) {
		return make([]interface{}, 0)
	}

	var dayOfWeek int64
	if input.DayOfWeek != nil {
		dayOfWeek = *input.DayOfWeek
	}
	var startHour int64
	if input.StartHour != nil {
		startHour = *input.StartHour
	}
	var startMinute int64
	if input.StartMinute != nil {
		startMinute = *input.StartMinute
	}
	return []interface{}{
		map[string]interface{}{
			"day_of_week":  dayOfWeek,
			"start_hour":   startHour,
			"start_minute": startMinute,
		},
	}
}

func flattenDataSourceFlexibleServerHighAvailability(ha *servers.HighAvailability) []interface{} {
	if ha == nil || *ha.Mode == servers.HighAvailabilityModeDisabled {
		return []interface{}{}
	}

	var zone string
	if ha.StandbyAvailabilityZone != nil {
		zone = *ha.StandbyAvailabilityZone
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                      string(*ha.Mode),
			"standby_availability_zone": zone,
		},
	}
}
