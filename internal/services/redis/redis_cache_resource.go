// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/patchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/redis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var skuWeight = map[string]int8{
	"Basic":    1,
	"Standard": 2,
	"Premium":  3,
}

func resourceRedisCache() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRedisCacheCreate,
		Read:   resourceRedisCacheRead,
		Update: resourceRedisCacheUpdate,
		Delete: resourceRedisCacheDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := redis.ParseRediID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.RedisCacheV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"capacity": {
				Type:     pluginsdk.TypeInt,
				Required: true,
			},

			"family": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// TODO: this should become case-sensitive in v4.0 (true -> false and remove DiffSupressFunc)
				ValidateFunc:     validation.StringInSlice(redis.PossibleValuesForSkuFamily(), true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.SkuNameBasic),
					string(redis.SkuNameStandard),
					string(redis.SkuNamePremium),
				}, false),
			},

			"minimum_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(redis.TlsVersionOnePointTwo),
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.TlsVersionOnePointZero),
					string(redis.TlsVersionOnePointOne),
					string(redis.TlsVersionOnePointTwo),
				}, false),
			},

			"shard_count": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_non_ssl_port": {
				Type:     pluginsdk.TypeBool,
				Default:  false,
				Optional: true,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"private_static_ip_address": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"redis_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"maxclients": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"maxmemory_delta": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Computed: true,
						},

						"maxmemory_reserved": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Computed: true,
						},

						"maxmemory_policy": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "volatile-lru",
							ValidateFunc: validate.MaxMemoryPolicy,
						},

						"maxfragmentationmemory_reserved": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Computed: true,
						},

						"rdb_backup_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"rdb_backup_frequency": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validate.CacheBackupFrequency,
						},

						"rdb_backup_max_snapshot_count": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},

						"rdb_storage_connection_string": {
							Type:      pluginsdk.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"notify_keyspace_events": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"aof_backup_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"aof_storage_connection_string_0": {
							Type:      pluginsdk.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"aof_storage_connection_string_1": {
							Type:      pluginsdk.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						// TODO 4.0: change this from enable_* to *_enabled
						"enable_authentication": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"patch_schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},

						"maintenance_window": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "PT5H",
							ValidateFunc: azValidate.ISO8601Duration,
						},

						"start_hour_utc": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
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

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// todo: investigate the difference between `replicas_per_master` and `replicas_per_primary` - are these
			// the same field that's been renamed ala Redis? https://github.com/Azure/azure-rest-api-specs/pull/13005
			"replicas_per_master": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
				// Can't make more than 3 replicas in portal, assuming it's a limitation
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"replicas_per_primary": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},

			"tenant_settings": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"redis_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"4", "6"}, false),
				DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
					n := strings.Split(old, ".")
					if len(n) >= 1 {
						newMajor := n[0]
						return new == newMajor
					}
					return false
				},
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, meta interface{}) bool {
				// downgrade the SKU is not supported, recreate the resource
				if old.(string) != "" && new.(string) != "" {
					return skuWeight[old.(string)] > skuWeight[new.(string)]
				}
				return false
			}),
		),
	}
}

func resourceRedisCacheCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	patchClient := meta.(*clients.Client).Redis.PatchSchedules
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := redis.NewRediID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_redis_cache", id.ID())
	}

	patchSchedule := expandRedisPatchSchedule(d)
	redisConfiguration, err := expandRedisConfiguration(d)
	if err != nil {
		return fmt.Errorf("parsing Redis Configuration: %+v", err)
	}

	publicNetworkAccess := redis.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = redis.PublicNetworkAccessDisabled
	}

	redisIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	parameters := redis.RedisCreateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: redis.RedisCreateProperties{
			EnableNonSslPort: pointer.To(d.Get("enable_non_ssl_port").(bool)),
			Sku: redis.Sku{
				Capacity: int64(d.Get("capacity").(int)),
				Family:   redis.SkuFamily(d.Get("family").(string)),
				Name:     redis.SkuName(d.Get("sku_name").(string)),
			},
			MinimumTlsVersion:   pointer.To(redis.TlsVersion(d.Get("minimum_tls_version").(string))),
			RedisConfiguration:  redisConfiguration,
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
		},
		Identity: redisIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("shard_count"); ok {
		shardCount := int64(v.(int))
		parameters.Properties.ShardCount = &shardCount
	}

	if v, ok := d.GetOk("replicas_per_master"); ok {
		parameters.Properties.ReplicasPerMaster = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("replicas_per_primary"); ok {
		parameters.Properties.ReplicasPerPrimary = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("redis_version"); ok {
		parameters.Properties.RedisVersion = utils.String(v.(string))
	}

	if v, ok := d.GetOk("tenant_settings"); ok {
		parameters.Properties.TenantSettings = expandTenantSettings(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("private_static_ip_address"); ok {
		parameters.Properties.StaticIP = utils.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		parsed, parseErr := commonids.ParseSubnetID(v.(string))
		if parseErr != nil {
			return err
		}

		locks.ByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(parsed.SubnetName, network.SubnetResourceName)
		defer locks.UnlockByName(parsed.SubnetName, network.SubnetResourceName)

		parameters.Properties.SubnetId = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := zones.ExpandUntyped(v.(*schema.Set).List())
		if len(zones) > 0 {
			parameters.Zones = &zones
		}
	}

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to become available", id)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(deadline),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	d.SetId(id.ID())

	if patchSchedule != nil {
		patchScheduleRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
		if _, err = patchClient.CreateOrUpdate(ctx, patchScheduleRedisId, *patchSchedule); err != nil {
			return fmt.Errorf("setting Patch Schedule for %s: %+v", id, err)
		}
	}

	return resourceRedisCacheRead(d, meta)
}

func resourceRedisCacheUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	patchClient := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redis.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := redis.RedisUpdateParameters{
		Properties: &redis.RedisUpdateProperties{
			MinimumTlsVersion: pointer.To(redis.TlsVersion(d.Get("minimum_tls_version").(string))),
			EnableNonSslPort:  utils.Bool(enableNonSSLPort),
			Sku: &redis.Sku{
				Capacity: int64(d.Get("capacity").(int)),
				Family:   redis.SkuFamily(d.Get("family").(string)),
				Name:     redis.SkuName(d.Get("sku_name").(string)),
			},
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("shard_count"); ok {
		if d.HasChange("shard_count") {
			shardCount := int64(v.(int))
			parameters.Properties.ShardCount = &shardCount
		}
	}

	if v, ok := d.GetOk("replicas_per_master"); ok {
		if d.HasChange("replicas_per_master") {
			parameters.Properties.ReplicasPerMaster = utils.Int64(int64(v.(int)))
		}
	}

	if v, ok := d.GetOk("replicas_per_primary"); ok {
		if d.HasChange("replicas_per_primary") {
			parameters.Properties.ReplicasPerPrimary = utils.Int64(int64(v.(int)))
		}
	}

	if v, ok := d.GetOk("redis_version"); ok {
		if d.HasChange("redis_version") {
			parameters.Properties.RedisVersion = utils.String(v.(string))
		}
	}

	if v, ok := d.GetOk("tenant_settings"); ok {
		if d.HasChange("tenant_settings") {
			parameters.Properties.TenantSettings = expandTenantSettings(v.(map[string]interface{}))
		}
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := redis.PublicNetworkAccessEnabled
		if !d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = redis.PublicNetworkAccessDisabled
		}
		parameters.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("redis_configuration") {
		redisConfiguration, err := expandRedisConfiguration(d)
		if err != nil {
			return fmt.Errorf("parsing Redis Configuration: %+v", err)
		}
		parameters.Properties.RedisConfiguration = redisConfiguration
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %s to become available", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating", "UpgradingRedisServerVersion"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	// identity cannot be updated with sku,publicNetworkAccess,redisVersion etc.
	if d.HasChange("identity") {
		redisIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf(`expanding "identity": %v`, err)
		}

		identityParameter := redis.RedisUpdateParameters{
			Identity: redisIdentity,
		}
		if _, err := client.Update(ctx, *id, identityParameter); err != nil {
			return fmt.Errorf("updating identity for %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Waiting for %s to become available", id)
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for %s to become available: %+v", id, err)
		}
	}

	if d.HasChange("patch_schedule") {
		patchSchedule := expandRedisPatchSchedule(d)

		patchSchedulesRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
		if patchSchedule == nil || len(patchSchedule.Properties.ScheduleEntries) == 0 {
			_, err = patchClient.Delete(ctx, patchSchedulesRedisId)
			if err != nil {
				return fmt.Errorf("deleting Patch Schedule for %s: %+v", *id, err)
			}
		} else {
			_, err = patchClient.CreateOrUpdate(ctx, patchSchedulesRedisId, *patchSchedule)
			if err != nil {
				return fmt.Errorf("setting Patch Schedule for %s: %+v", *id, err)
			}
		}
	}

	return resourceRedisCacheRead(d, meta)
}

func resourceRedisCacheRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	patchSchedulesClient := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redis.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	patchSchedulesRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
	schedule, err := patchSchedulesClient.Get(ctx, patchSchedulesRedisId)
	var patchSchedule []interface{}
	if err == nil {
		patchSchedule = flattenRedisPatchSchedules(*schedule.Model)
	}
	if err = d.Set("patch_schedule", patchSchedule); err != nil {
		return fmt.Errorf("setting `patch_schedule`: %+v", err)
	}

	d.Set("name", id.RedisName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		redisIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", redisIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		d.Set("location", location.Normalize(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		props := model.Properties
		d.Set("capacity", int(props.Sku.Capacity))
		d.Set("family", string(props.Sku.Family))
		d.Set("sku_name", string(props.Sku.Name))

		d.Set("ssl_port", props.SslPort)
		d.Set("hostname", props.HostName)
		minimumTlsVersion := string(redis.TlsVersionOnePointTwo)
		if props.MinimumTlsVersion != nil {
			minimumTlsVersion = string(*props.MinimumTlsVersion)
		}
		d.Set("minimum_tls_version", minimumTlsVersion)
		d.Set("port", props.Port)
		d.Set("enable_non_ssl_port", props.EnableNonSslPort)
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

		publicNetworkAccessEnabled := true
		if props.PublicNetworkAccess != nil {
			publicNetworkAccessEnabled = *props.PublicNetworkAccess == redis.PublicNetworkAccessEnabled
		}
		d.Set("public_network_access_enabled", publicNetworkAccessEnabled)
		d.Set("replicas_per_master", props.ReplicasPerMaster)
		d.Set("replicas_per_primary", props.ReplicasPerPrimary)
		d.Set("redis_version", props.RedisVersion)
		d.Set("tenant_settings", flattenTenantSettings(props.TenantSettings))

		redisConfiguration, err := flattenRedisConfiguration(props.RedisConfiguration)
		if err != nil {
			return fmt.Errorf("flattening `redis_configuration`: %+v", err)
		}
		if err := d.Set("redis_configuration", redisConfiguration); err != nil {
			return fmt.Errorf("setting `redis_configuration`: %+v", err)
		}

		enableSslPort := !*props.EnableNonSslPort
		d.Set("primary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.Model.PrimaryKey, enableSslPort))
		d.Set("secondary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.Model.SecondaryKey, enableSslPort))
		d.Set("primary_access_key", keysResp.Model.PrimaryKey)
		d.Set("secondary_access_key", keysResp.Model.SecondaryKey)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceRedisCacheDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redis.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if subnetID := read.Model.Properties.SubnetId; subnetID != nil {
		parsed, parseErr := commonids.ParseSubnetIDInsensitively(*subnetID)
		if parseErr != nil {
			return err
		}

		locks.ByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(parsed.SubnetName, network.SubnetResourceName)
		defer locks.UnlockByName(parsed.SubnetName, network.SubnetResourceName)
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func redisStateRefreshFunc(ctx context.Context, client *redis.RedisClient, id redis.RediId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("polling for status of %s: %+v", id, err)
		}

		provisioningState := ""
		if model := res.Model; model != nil && model.Properties.ProvisioningState != nil {
			provisioningState = string(*res.Model.Properties.ProvisioningState)
		}
		if provisioningState == "" {
			return nil, "", fmt.Errorf("polling for status of %s: `provisioningState` was nil", id)
		}

		return res, provisioningState, nil
	}
}

func expandRedisConfiguration(d *pluginsdk.ResourceData) (*redis.RedisCommonPropertiesRedisConfiguration, error) {
	output := &redis.RedisCommonPropertiesRedisConfiguration{}

	input := d.Get("redis_configuration").([]interface{})
	if len(input) == 0 || input[0] == nil {
		return output, nil
	}
	raw := input[0].(map[string]interface{})
	skuName := d.Get("sku_name").(string)

	if v := raw["maxclients"].(int); v > 0 {
		output.Maxclients = utils.String(strconv.Itoa(v))
	}

	if d.Get("sku_name").(string) != string(redis.SkuNameBasic) {
		if v := raw["maxmemory_delta"].(int); v > 0 {
			output.MaxmemoryDelta = utils.String(strconv.Itoa(v))
		}

		if v := raw["maxmemory_reserved"].(int); v > 0 {
			output.MaxmemoryReserved = utils.String(strconv.Itoa(v))
		}

		if v := raw["maxfragmentationmemory_reserved"].(int); v > 0 {
			output.MaxfragmentationmemoryReserved = utils.String(strconv.Itoa(v))
		}
	}

	if v := raw["maxmemory_policy"].(string); v != "" {
		output.MaxmemoryPolicy = utils.String(v)
	}

	// RDB Backup
	// nolint : staticcheck
	v, valExists := d.GetOkExists("redis_configuration.0.rdb_backup_enabled")
	if valExists {
		rdbBackupEnabled := v.(bool)

		// rdb_backup_enabled is available when SKU is Premium
		if strings.EqualFold(skuName, string(redis.SkuNamePremium)) {
			if rdbBackupEnabled {
				if connStr := raw["rdb_storage_connection_string"].(string); connStr == "" {
					return nil, fmt.Errorf("The rdb_storage_connection_string property must be set when rdb_backup_enabled is true")
				}
			}
			output.RdbBackupEnabled = utils.String(strconv.FormatBool(rdbBackupEnabled))
		} else if rdbBackupEnabled && !strings.EqualFold(skuName, string(redis.SkuNamePremium)) {
			return nil, fmt.Errorf("The `rdb_backup_enabled` property requires a `Premium` sku to be set")
		}
	}

	if v := raw["rdb_backup_frequency"].(int); v > 0 {
		output.RdbBackupFrequency = utils.String(strconv.Itoa(v))
	}

	if v := raw["rdb_backup_max_snapshot_count"].(int); v > 0 {
		output.RdbBackupMaxSnapshotCount = utils.String(strconv.Itoa(v))
	}

	if v := raw["rdb_storage_connection_string"].(string); v != "" {
		output.RdbStorageConnectionString = utils.String(v)
	}

	if v := raw["notify_keyspace_events"].(string); v != "" {
		output.NotifyKeyspaceEvents = utils.String(v)
	}

	// AOF Backup
	// nolint : staticcheck
	v, valExists = d.GetOkExists("redis_configuration.0.aof_backup_enabled")
	if valExists {
		// aof_backup_enabled is available when SKU is Premium
		if strings.EqualFold(skuName, string(redis.SkuNamePremium)) {
			output.AofBackupEnabled = utils.String(strconv.FormatBool(v.(bool)))
		}
	}

	if v := raw["aof_storage_connection_string_0"].(string); v != "" {
		output.AofStorageConnectionString0 = utils.String(v)
	}

	if v := raw["aof_storage_connection_string_1"].(string); v != "" {
		output.AofStorageConnectionString1 = utils.String(v)
	}
	authEnabled := raw["enable_authentication"].(bool)
	// Redis authentication can only be disabled if it is launched inside a VNET.
	if _, isPrivate := d.GetOk("subnet_id"); !isPrivate {
		if !authEnabled {
			return nil, fmt.Errorf("Cannot set `enable_authentication` to `false` when `subnet_id` is not set")
		}
	} else {
		value := isAuthNotRequiredAsString(authEnabled)
		output.Authnotrequired = utils.String(value)
	}
	return output, nil
}

func expandRedisPatchSchedule(d *pluginsdk.ResourceData) *patchschedules.RedisPatchSchedule {
	v, ok := d.GetOk("patch_schedule")
	if !ok {
		return nil
	}

	scheduleValues := v.([]interface{})
	entries := make([]patchschedules.ScheduleEntry, 0)
	for _, scheduleValue := range scheduleValues {
		vals := scheduleValue.(map[string]interface{})
		dayOfWeek := vals["day_of_week"].(string)
		maintenanceWindow := vals["maintenance_window"].(string)
		startHourUtc := vals["start_hour_utc"].(int)

		entries = append(entries, patchschedules.ScheduleEntry{
			DayOfWeek:         patchschedules.DayOfWeek(dayOfWeek),
			MaintenanceWindow: utils.String(maintenanceWindow),
			StartHourUtc:      int64(startHourUtc),
		})
	}

	schedule := patchschedules.RedisPatchSchedule{
		Properties: patchschedules.ScheduleEntries{
			ScheduleEntries: entries,
		},
	}
	return &schedule
}

func expandTenantSettings(input map[string]interface{}) *map[string]string {
	output := make(map[string]string, len(input))

	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}

func flattenTenantSettings(input *map[string]string) map[string]string {
	output := make(map[string]string)

	if input != nil {
		for k, v := range *input {
			output[k] = v
		}
	}

	return output
}

func flattenRedisConfiguration(input *redis.RedisCommonPropertiesRedisConfiguration) ([]interface{}, error) {
	outputs := make(map[string]interface{})

	if input.Maxclients != nil {
		i, err := strconv.Atoi(*input.Maxclients)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxclients` %q: %+v", *input.Maxclients, err)
		}
		outputs["maxclients"] = i
	}
	if input.MaxmemoryDelta != nil {
		i, err := strconv.Atoi(*input.MaxmemoryDelta)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxmemory-delta` %q: %+v", *input.MaxmemoryDelta, err)
		}
		outputs["maxmemory_delta"] = i
	}
	if input.MaxmemoryReserved != nil {
		i, err := strconv.Atoi(*input.MaxmemoryReserved)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxmemory-reserved` %q: %+v", *input.MaxmemoryReserved, err)
		}
		outputs["maxmemory_reserved"] = i
	}
	if input.MaxmemoryPolicy != nil {
		outputs["maxmemory_policy"] = *input.MaxmemoryPolicy
	}

	if input.MaxfragmentationmemoryReserved != nil {
		i, err := strconv.Atoi(*input.MaxfragmentationmemoryReserved)
		if err != nil {
			return nil, fmt.Errorf("parsing `maxfragmentationmemory-reserved` %q: %+v", *input.MaxfragmentationmemoryReserved, err)
		}
		outputs["maxfragmentationmemory_reserved"] = i
	}

	// delta, reserved, enabled, frequency,, count,
	if input.RdbBackupEnabled != nil {
		b, err := strconv.ParseBool(*input.RdbBackupEnabled)
		if err != nil {
			return nil, fmt.Errorf("parsing `rdb-backup-enabled` %q: %+v", *input.RdbBackupEnabled, err)
		}
		outputs["rdb_backup_enabled"] = b
	}
	if input.RdbBackupFrequency != nil {
		i, err := strconv.Atoi(*input.RdbBackupFrequency)
		if err != nil {
			return nil, fmt.Errorf("parsing `rdb-backup-frequency` %q: %+v", *input.RdbBackupFrequency, err)
		}
		outputs["rdb_backup_frequency"] = i
	}
	if input.RdbBackupMaxSnapshotCount != nil {
		i, err := strconv.Atoi(*input.RdbBackupMaxSnapshotCount)
		if err != nil {
			return nil, fmt.Errorf("parsing `rdb-backup-max-snapshot-count` %q: %+v", *input.RdbBackupMaxSnapshotCount, err)
		}
		outputs["rdb_backup_max_snapshot_count"] = i
	}
	if input.RdbStorageConnectionString != nil {
		outputs["rdb_storage_connection_string"] = *input.RdbStorageConnectionString
	}
	if v := input.NotifyKeyspaceEvents; v != nil {
		outputs["notify_keyspace_events"] = v
	}

	if v := input.AofBackupEnabled; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("parsing `aof-backup-enabled` %q: %+v", *v, err)
		}
		outputs["aof_backup_enabled"] = b
	}
	if input.AofStorageConnectionString0 != nil {
		outputs["aof_storage_connection_string_0"] = *input.AofStorageConnectionString0
	}
	if input.AofStorageConnectionString1 != nil {
		outputs["aof_storage_connection_string_1"] = *input.AofStorageConnectionString1
	}

	// `authnotrequired` is not set for instances launched outside a VNET
	outputs["enable_authentication"] = true
	if v := input.Authnotrequired; v != nil {
		outputs["enable_authentication"] = isAuthRequiredAsBool(*v)
	}

	return []interface{}{outputs}, nil
}

func isAuthRequiredAsBool(notRequired string) bool {
	value := strings.ToLower(notRequired)
	output := map[string]bool{
		"yes": false,
		"no":  true,
	}
	return output[value]
}

func isAuthNotRequiredAsString(authRequired bool) string {
	output := map[bool]string{
		true:  "no",
		false: "yes",
	}
	return output[authRequired]
}

func flattenRedisPatchSchedules(schedule patchschedules.RedisPatchSchedule) []interface{} {
	outputs := make([]interface{}, 0)

	for _, entry := range schedule.Properties.ScheduleEntries {
		maintenanceWindow := ""
		if entry.MaintenanceWindow != nil {
			maintenanceWindow = *entry.MaintenanceWindow
		}

		outputs = append(outputs, map[string]interface{}{
			"day_of_week":        string(entry.DayOfWeek),
			"maintenance_window": maintenanceWindow,
			"start_hour_utc":     int(entry.StartHourUtc),
		})
	}

	return outputs
}

func getRedisConnectionString(redisHostName string, sslPort int64, accessKey string, enableSslPort bool) string {
	return fmt.Sprintf("%s:%d,password=%s,ssl=%t,abortConnect=False", redisHostName, sslPort, accessKey, enableSslPort)
}
