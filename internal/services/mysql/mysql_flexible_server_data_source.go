package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2021-05-01/mysqlflexibleservers"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMysqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFlexibleServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("zone", props.AvailabilityZone)
		d.Set("version", props.Version)
		d.Set("fqdn", props.FullyQualifiedDomainName)

		if network := props.Network; network != nil {
			d.Set("public_network_access_enabled", network.PublicNetworkAccess == mysqlflexibleservers.EnableStatusEnumEnabled)
			d.Set("delegated_subnet_id", network.DelegatedSubnetResourceID)
			d.Set("private_dns_zone_id", network.PrivateDNSZoneResourceID)
		}

		if err := d.Set("maintenance_window", flattenDataSourceArmServerMaintenanceWindow(props.MaintenanceWindow)); err != nil {
			return fmt.Errorf("setting `maintenance_window`: %+v", err)
		}

		if err := d.Set("storage", flattenDataSourceArmServerStorage(props.Storage)); err != nil {
			return fmt.Errorf("setting `storage`: %+v", err)
		}

		if backup := props.Backup; backup != nil {
			d.Set("backup_retention_days", backup.BackupRetentionDays)
			d.Set("geo_redundant_backup_enabled", backup.GeoRedundantBackup == mysqlflexibleservers.EnableStatusEnumEnabled)
		}

		if err := d.Set("high_availability", flattenDataSourceFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
			return fmt.Errorf("setting `high_availability`: %+v", err)
		}
		d.Set("replication_role", props.ReplicationRole)
		d.Set("replica_capacity", props.ReplicaCapacity)
	}

	sku, err := flattenDataSourceFlexibleServerSku(resp.Sku)
	if err != nil {
		return fmt.Errorf("flattening `sku_name` for %q: %v", id, err)
	}

	d.Set("sku_name", sku)

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenDataSourceArmServerStorage(storage *mysqlflexibleservers.Storage) []interface{} {
	if storage == nil {
		return []interface{}{}
	}

	var size, iops int32
	if storage.StorageSizeGB != nil {
		size = *storage.StorageSizeGB
	}

	if storage.Iops != nil {
		iops = *storage.Iops
	}

	return []interface{}{
		map[string]interface{}{
			"size_gb":           size,
			"iops":              iops,
			"auto_grow_enabled": storage.AutoGrow == mysqlflexibleservers.EnableStatusEnumEnabled,
		},
	}
}

func flattenDataSourceFlexibleServerSku(sku *mysqlflexibleservers.Sku) (string, error) {
	if sku == nil || sku.Name == nil || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case mysqlflexibleservers.SkuTierBurstable:
		tier = "B"
	case mysqlflexibleservers.SkuTierGeneralPurpose:
		tier = "GP"
	case mysqlflexibleservers.SkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, *sku.Name}, "_"), nil
}

func flattenDataSourceArmServerMaintenanceWindow(input *mysqlflexibleservers.MaintenanceWindow) []interface{} {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == string(ServerMaintenanceWindowDisabled) {
		return make([]interface{}, 0)
	}

	var dayOfWeek int32
	if input.DayOfWeek != nil {
		dayOfWeek = *input.DayOfWeek
	}
	var startHour int32
	if input.StartHour != nil {
		startHour = *input.StartHour
	}
	var startMinute int32
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

func flattenDataSourceFlexibleServerHighAvailability(ha *mysqlflexibleservers.HighAvailability) []interface{} {
	if ha == nil || ha.Mode == mysqlflexibleservers.HighAvailabilityModeDisabled {
		return []interface{}{}
	}

	var zone string
	if ha.StandbyAvailabilityZone != nil {
		zone = *ha.StandbyAvailabilityZone
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                      string(ha.Mode),
			"standby_availability_zone": zone,
		},
	}
}
