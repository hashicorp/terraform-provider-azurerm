// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/patchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceRedisCache() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceRedisCacheRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"zones": commonschema.ZonesMultipleComputed(),

			"capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"family": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"minimum_tls_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"shard_count": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"non_ssl_port_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_static_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"redis_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"active_directory_authentication_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"maxclients": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"maxmemory_delta": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"maxmemory_reserved": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"maxmemory_policy": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"maxfragmentationmemory_reserved": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"rdb_backup_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"rdb_backup_frequency": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"rdb_backup_max_snapshot_count": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"rdb_storage_connection_string": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"notify_keyspace_events": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"aof_backup_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"aof_storage_connection_string_0": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"aof_storage_connection_string_1": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"authentication_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"storage_account_subscription_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"data_persistence_authentication_method": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"patch_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"maintenance_window": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"start_hour_utc": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"ssl_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"access_keys_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["enable_non_ssl_port"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "`enable_non_ssl_port` will be removed in favour of the property `non_ssl_port_enabled` in version 4.0 of the AzureRM Provider.",
		}
		resource.Schema["redis_configuration"].Elem.(*pluginsdk.Resource).Schema["enable_authentication"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Default:    true,
			Deprecated: "`enable_authentication` will be removed in favour of the property `authentication_enabled` in version 4.0 of the AzureRM Provider.",
		}
	}

	return resource
}

func dataSourceRedisCacheRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	patchSchedulesClient := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := redis.NewRediID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	patchScheduleRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
	schedule, err := patchSchedulesClient.Get(ctx, patchScheduleRedisId)
	if err != nil {
		if !response.WasNotFound(schedule.HttpResponse) {
			return fmt.Errorf("obtaining patch schedules for %s: %+v", id, err)
		}
	}
	var patchSchedule []interface{}
	if model := schedule.Model; model != nil {
		patchSchedule = flattenRedisPatchSchedules(*schedule.Model)
	}

	keys, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		props := model.Properties
		sku := props.Sku
		d.Set("capacity", sku.Capacity)
		d.Set("family", sku.Family)
		d.Set("sku_name", sku.Name)

		d.Set("ssl_port", props.SslPort)
		d.Set("hostname", props.HostName)
		minimumTlsVersion := string(redis.TlsVersionOnePointTwo)
		if props.MinimumTlsVersion != nil {
			minimumTlsVersion = string(*props.MinimumTlsVersion)
		}
		d.Set("minimum_tls_version", minimumTlsVersion)
		d.Set("port", props.Port)
		d.Set("non_ssl_port_enabled", props.EnableNonSslPort)

		if !features.FourPointOhBeta() {
			d.Set("enable_non_ssl_port", props.EnableNonSslPort)
		}

		shardCount := 0
		if props.ShardCount != nil {
			shardCount = int(*props.ShardCount)
		}
		d.Set("shard_count", shardCount)
		d.Set("private_static_ip_address", props.StaticIP)
		subnetId := ""
		if props.SubnetId != nil {
			parsed, err := commonids.ParseSubnetIDInsensitively(*props.SubnetId)
			if err != nil {
				return err
			}

			subnetId = parsed.ID()
		}
		d.Set("subnet_id", subnetId)

		redisConfiguration, err := flattenRedisConfiguration(props.RedisConfiguration)
		if err != nil {
			return fmt.Errorf("flattening `redis_configuration`: %+v", err)
		}
		if err := d.Set("redis_configuration", redisConfiguration); err != nil {
			return fmt.Errorf("setting `redis_configuration`: %+v", err)
		}

		enableSslPort := !*props.EnableNonSslPort
		d.Set("primary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keys.Model.PrimaryKey, enableSslPort))
		d.Set("secondary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keys.Model.SecondaryKey, enableSslPort))
		d.Set("access_keys_authentication_enabled", !pointer.From(props.DisableAccessKeyAuthentication))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if err = d.Set("patch_schedule", patchSchedule); err != nil {
		return fmt.Errorf("setting `patch_schedule`: %+v", err)
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}
