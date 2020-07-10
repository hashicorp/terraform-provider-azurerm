package redis

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRedisCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRedisCacheCreate,
		Read:   resourceArmRedisCacheRead,
		Update: resourceArmRedisCacheUpdate,
		Delete: resourceArmRedisCacheDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"zones": azure.SchemaSingleZone(),

			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"family": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateRedisFamily,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.Basic),
					string(redis.Standard),
					string(redis.Premium),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"minimum_tls_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  redis.OneFullStopZero,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.OneFullStopZero),
					string(redis.OneFullStopOne),
					string(redis.OneFullStopTwo),
				}, false),
			},

			"shard_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"enable_non_ssl_port": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"private_static_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"redis_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maxclients": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"maxmemory_delta": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"maxmemory_reserved": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"maxmemory_policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "volatile-lru",
							ValidateFunc: validateRedisMaxMemoryPolicy,
						},

						"maxfragmentationmemory_reserved": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"rdb_backup_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"rdb_backup_frequency": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateRedisBackupFrequency,
						},

						"rdb_backup_max_snapshot_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"rdb_storage_connection_string": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"notify_keyspace_events": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"aof_backup_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"aof_storage_connection_string_0": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},

						"aof_storage_connection_string_1": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"enable_authentication": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"patch_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},
						"start_hour_utc": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
				},
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ssl_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmRedisCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Redis Instance %s (resource group %s) ID", name, resGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_redis_cache", *existing.ID)
		}
	}

	patchSchedule := expandRedisPatchSchedule(d)
	redisConfiguration, err := expandRedisConfiguration(d)
	if err != nil {
		return fmt.Errorf("Error parsing Redis Configuration: %+v", err)
	}

	parameters := redis.CreateParameters{
		Location: utils.String(location),
		CreateProperties: &redis.CreateProperties{
			EnableNonSslPort: utils.Bool(enableNonSSLPort),
			Sku: &redis.Sku{
				Capacity: utils.Int32(capacity),
				Family:   family,
				Name:     sku,
			},
			MinimumTLSVersion:  redis.TLSVersion(d.Get("minimum_tls_version").(string)),
			RedisConfiguration: redisConfiguration,
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("shard_count"); ok {
		shardCount := int32(v.(int))
		parameters.ShardCount = &shardCount
	}

	if v, ok := d.GetOk("private_static_ip_address"); ok {
		parameters.StaticIP = utils.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		parsed, parseErr := azure.ParseAzureResourceID(v.(string))
		if parseErr != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", v.(string))
		}
		subnetName := parsed.Path["subnets"]
		virtualNetworkName := parsed.Path["virtualNetworks"]

		locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(subnetName, network.SubnetResourceName)
		defer locks.UnlockByName(subnetName, network.SubnetResourceName)

		parameters.SubnetID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		parameters.Zones = azure.ExpandZones(v.([]interface{}))
	}

	future, err := client.Create(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create request for Redis Cache %s (resource group %s): %v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the create of Redis Cache %s (resource group %s): %v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading Redis Cache %s (resource group %s): %v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Redis Cache %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for Redis Cache (%s) to become available", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, resGroup, name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Cache (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	if schedule := patchSchedule; schedule != nil {
		patchClient := meta.(*clients.Client).Redis.PatchSchedulesClient
		_, err = patchClient.CreateOrUpdate(ctx, resGroup, name, *schedule)
		if err != nil {
			return fmt.Errorf("Error setting Redis Patch Schedule: %+v", err)
		}
	}

	return resourceArmRedisCacheRead(d, meta)
}

func resourceArmRedisCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache update.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := redis.UpdateParameters{
		UpdateProperties: &redis.UpdateProperties{
			MinimumTLSVersion: redis.TLSVersion(d.Get("minimum_tls_version").(string)),
			EnableNonSslPort:  utils.Bool(enableNonSSLPort),
			Sku: &redis.Sku{
				Capacity: utils.Int32(capacity),
				Family:   family,
				Name:     sku,
			},
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("shard_count"); ok {
		if d.HasChange("shard_count") {
			shardCount := int32(v.(int))
			parameters.ShardCount = &shardCount
		}
	}

	if d.HasChange("redis_configuration") {
		redisConfiguration, err := expandRedisConfiguration(d)
		if err != nil {
			return fmt.Errorf("Error parsing Redis Configuration: %+v", err)
		}
		parameters.RedisConfiguration = redisConfiguration
	}

	if _, err := client.Update(ctx, resGroup, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Redis Instance %s (resource group %s) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for Redis Instance (%s) to become available", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, resGroup, name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Instance (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	patchSchedule := expandRedisPatchSchedule(d)

	patchClient := meta.(*clients.Client).Redis.PatchSchedulesClient
	if patchSchedule == nil || len(*patchSchedule.ScheduleEntries.ScheduleEntries) == 0 {
		_, err = patchClient.Delete(ctx, resGroup, name)
		if err != nil {
			return fmt.Errorf("Error deleting Redis Patch Schedule: %+v", err)
		}
	} else {
		_, err = patchClient.CreateOrUpdate(ctx, resGroup, name, *patchSchedule)
		if err != nil {
			return fmt.Errorf("Error setting Redis Patch Schedule: %+v", err)
		}
	}

	return resourceArmRedisCacheRead(d, meta)
}

func resourceArmRedisCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	resp, err := client.Get(ctx, resGroup, name)

	// covers if the resource has been deleted outside of TF, but is still in the state
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Redis Cache %s: %s", name, err)
	}

	keysResp, err := client.ListKeys(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making ListKeys request on Azure Redis Cache %s: %s", name, err)
	}

	patchSchedulesClient := meta.(*clients.Client).Redis.PatchSchedulesClient

	schedule, err := patchSchedulesClient.Get(ctx, resGroup, name)
	if err == nil {
		patchSchedule := flattenRedisPatchSchedules(schedule)
		if err = d.Set("patch_schedule", patchSchedule); err != nil {
			return fmt.Errorf("Error setting `patch_schedule`: %+v", err)
		}
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if zones := resp.Zones; zones != nil {
		d.Set("zones", zones)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("capacity", sku.Capacity)
		d.Set("family", sku.Family)
		d.Set("sku_name", sku.Name)
	}

	props := resp.Properties
	if props != nil {
		d.Set("ssl_port", props.SslPort)
		d.Set("hostname", props.HostName)
		d.Set("minimum_tls_version", string(props.MinimumTLSVersion))
		d.Set("port", props.Port)
		d.Set("enable_non_ssl_port", props.EnableNonSslPort)
		if props.ShardCount != nil {
			d.Set("shard_count", props.ShardCount)
		}
		d.Set("private_static_ip_address", props.StaticIP)
		d.Set("subnet_id", props.SubnetID)
	}

	redisConfiguration, err := flattenRedisConfiguration(resp.RedisConfiguration)
	if err != nil {
		return fmt.Errorf("Error flattening `redis_configuration`: %+v", err)
	}
	if err := d.Set("redis_configuration", redisConfiguration); err != nil {
		return fmt.Errorf("Error setting `redis_configuration`: %+v", err)
	}

	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	if props != nil {
		enableSslPort := !*props.EnableNonSslPort
		d.Set("primary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.PrimaryKey, enableSslPort))
		d.Set("secondary_connection_string", getRedisConnectionString(*props.HostName, *props.SslPort, *keysResp.SecondaryKey, enableSslPort))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmRedisCacheDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Redis Cache %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.Properties == nil {
		return fmt.Errorf("Error retrieving Redis Cache properties %q (Resource Group %q): `props` was nil", name, resGroup)
	}
	props := *read.Properties
	if subnetID := props.SubnetID; subnetID != nil {
		parsed, parseErr := azure.ParseAzureResourceID(*subnetID)
		if parseErr != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", *subnetID)
		}
		subnetName := parsed.Path["subnets"]
		virtualNetworkName := parsed.Path["virtualNetworks"]

		locks.ByName(virtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(virtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(subnetName, network.SubnetResourceName)
		defer locks.UnlockByName(subnetName, network.SubnetResourceName)
	}
	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	return nil
}

func redisStateRefreshFunc(ctx context.Context, client *redis.Client, resourceGroupName string, sgName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, sgName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in redisStateRefreshFunc to Azure ARM for Redis Cache Instance '%s' (RG: '%s'): %s", sgName, resourceGroupName, err)
		}

		return res, string(res.ProvisioningState), nil
	}
}

func expandRedisConfiguration(d *schema.ResourceData) (map[string]*string, error) {
	output := make(map[string]*string)

	if v, ok := d.GetOk("redis_configuration.0.maxclients"); ok {
		clients := strconv.Itoa(v.(int))
		output["maxclients"] = utils.String(clients)
	}

	if v, ok := d.GetOk("redis_configuration.0.maxmemory_delta"); ok {
		delta := strconv.Itoa(v.(int))
		output["maxmemory-delta"] = utils.String(delta)
	}

	if v, ok := d.GetOk("redis_configuration.0.maxmemory_reserved"); ok {
		delta := strconv.Itoa(v.(int))
		output["maxmemory-reserved"] = utils.String(delta)
	}

	if v, ok := d.GetOk("redis_configuration.0.maxmemory_policy"); ok {
		output["maxmemory-policy"] = utils.String(v.(string))
	}

	if v, ok := d.GetOk("redis_configuration.0.maxfragmentationmemory_reserved"); ok {
		delta := strconv.Itoa(v.(int))
		output["maxfragmentationmemory-reserved"] = utils.String(delta)
	}

	// RDB Backup
	if v, ok := d.GetOk("redis_configuration.0.rdb_backup_enabled"); ok {
		if v.(bool) {
			if connStr, connOk := d.GetOk("redis_configuration.0.rdb_storage_connection_string"); !connOk || connStr.(string) == "" {
				return nil, fmt.Errorf("The rdb_storage_connection_string property must be set when rdb_backup_enabled is true")
			}
		}
		delta := strconv.FormatBool(v.(bool))
		output["rdb-backup-enabled"] = utils.String(delta)
	}

	if v, ok := d.GetOk("redis_configuration.0.rdb_backup_frequency"); ok {
		delta := strconv.Itoa(v.(int))
		output["rdb-backup-frequency"] = utils.String(delta)
	}

	if v, ok := d.GetOk("redis_configuration.0.rdb_backup_max_snapshot_count"); ok {
		delta := strconv.Itoa(v.(int))
		output["rdb-backup-max-snapshot-count"] = utils.String(delta)
	}

	if v, ok := d.GetOk("redis_configuration.0.rdb_storage_connection_string"); ok {
		output["rdb-storage-connection-string"] = utils.String(v.(string))
	}

	if v, ok := d.GetOk("redis_configuration.0.notify_keyspace_events"); ok {
		output["notify-keyspace-events"] = utils.String(v.(string))
	}

	// AOF Backup
	if v, ok := d.GetOk("redis_configuration.0.aof_backup_enabled"); ok {
		delta := strconv.FormatBool(v.(bool))
		output["aof-backup-enabled"] = utils.String(delta)
	}

	if v, ok := d.GetOk("redis_configuration.0.aof_storage_connection_string_0"); ok {
		output["aof-storage-connection-string-0"] = utils.String(v.(string))
	}

	if v, ok := d.GetOk("redis_configuration.0.aof_storage_connection_string_1"); ok {
		output["aof-storage-connection-string-1"] = utils.String(v.(string))
	}

	if authEnabled := d.Get("redis_configuration.0.enable_authentication").(bool); authEnabled {
		_, isPrivate := d.GetOk("subnet_id")

		// Redis authentication can only be disabled if it is launched inside a VNET.
		if !isPrivate {
			if !authEnabled {
				return nil, fmt.Errorf("Cannot set `enable_authentication` to `false` when `subnet_id` is not set")
			}
		} else {
			value := isAuthNotRequiredAsString(authEnabled)
			output["authnotrequired"] = utils.String(value)
		}
	}
	return output, nil
}

func expandRedisPatchSchedule(d *schema.ResourceData) *redis.PatchSchedule {
	v, ok := d.GetOk("patch_schedule")
	if !ok {
		return nil
	}

	scheduleValues := v.([]interface{})
	entries := make([]redis.ScheduleEntry, 0)
	for _, scheduleValue := range scheduleValues {
		vals := scheduleValue.(map[string]interface{})
		dayOfWeek := vals["day_of_week"].(string)
		startHourUtc := vals["start_hour_utc"].(int)

		entry := redis.ScheduleEntry{
			DayOfWeek:    redis.DayOfWeek(dayOfWeek),
			StartHourUtc: utils.Int32(int32(startHourUtc)),
		}
		entries = append(entries, entry)
	}

	schedule := redis.PatchSchedule{
		ScheduleEntries: &redis.ScheduleEntries{
			ScheduleEntries: &entries,
		},
	}
	return &schedule
}

func flattenRedisConfiguration(input map[string]*string) ([]interface{}, error) {
	outputs := make(map[string]interface{}, len(input))

	if v := input["maxclients"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxclients` %q: %+v", *v, err)
		}
		outputs["maxclients"] = i
	}
	if v := input["maxmemory-delta"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxmemory-delta` %q: %+v", *v, err)
		}
		outputs["maxmemory_delta"] = i
	}
	if v := input["maxmemory-reserved"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxmemory-reserved` %q: %+v", *v, err)
		}
		outputs["maxmemory_reserved"] = i
	}
	if v := input["maxmemory-policy"]; v != nil {
		outputs["maxmemory_policy"] = *v
	}

	if v := input["maxfragmentationmemory-reserved"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxfragmentationmemory-reserved` %q: %+v", *v, err)
		}
		outputs["maxfragmentationmemory_reserved"] = i
	}

	// delta, reserved, enabled, frequency,, count,
	if v := input["rdb-backup-enabled"]; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-enabled` %q: %+v", *v, err)
		}
		outputs["rdb_backup_enabled"] = b
	}
	if v := input["rdb-backup-frequency"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-frequency` %q: %+v", *v, err)
		}
		outputs["rdb_backup_frequency"] = i
	}
	if v := input["rdb-backup-max-snapshot-count"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-max-snapshot-count` %q: %+v", *v, err)
		}
		outputs["rdb_backup_max_snapshot_count"] = i
	}
	if v := input["rdb-storage-connection-string"]; v != nil {
		outputs["rdb_storage_connection_string"] = *v
	}
	if v := input["notify-keyspace-events"]; v != nil {
		outputs["notify_keyspace_events"] = *v
	}

	if v := input["aof-backup-enabled"]; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `aof-backup-enabled` %q: %+v", *v, err)
		}
		outputs["aof_backup_enabled"] = b
	}
	if v := input["aof-storage-connection-string-0"]; v != nil {
		outputs["aof_storage_connection_string_0"] = *v
	}
	if v := input["aof-storage-connection-string-1"]; v != nil {
		outputs["aof_storage_connection_string_1"] = *v
	}

	// `authnotrequired` is not set for instances launched outside a VNET
	outputs["enable_authentication"] = true
	if v := input["authnotrequired"]; v != nil {
		outputs["enable_authentication"] = isAuthRequiredAsBool(*v)
	}

	return []interface{}{outputs}, nil
}

func isAuthRequiredAsBool(not_required string) bool {
	value := strings.ToLower(not_required)
	output := map[string]bool{
		"yes": false,
		"no":  true,
	}
	return output[value]
}

func isAuthNotRequiredAsString(auth_required bool) string {
	output := map[bool]string{
		true:  "no",
		false: "yes",
	}
	return output[auth_required]
}

func flattenRedisPatchSchedules(schedule redis.PatchSchedule) []interface{} {
	outputs := make([]interface{}, 0)

	for _, entry := range *schedule.ScheduleEntries.ScheduleEntries {
		output := make(map[string]interface{})

		output["day_of_week"] = string(entry.DayOfWeek)
		output["start_hour_utc"] = int(*entry.StartHourUtc)

		outputs = append(outputs, output)
	}

	return outputs
}

func validateRedisFamily(v interface{}, _ string) (warnings []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"c": true,
		"p": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Family can only be C or P"))
	}
	return warnings, errors
}

func validateRedisMaxMemoryPolicy(v interface{}, _ string) (warnings []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"noeviction":      true,
		"allkeys-lru":     true,
		"volatile-lru":    true,
		"allkeys-random":  true,
		"volatile-random": true,
		"volatile-ttl":    true,
		"allkeys-lfu":     true,
		"volatile-lfu":    true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Max Memory Policy can only be 'noeviction' / 'allkeys-lru' / 'volatile-lru' / 'allkeys-random' / 'volatile-random' / 'volatile-ttl' / 'allkeys-lfu' / 'volatile-lfu'"))
	}

	return warnings, errors
}

func validateRedisBackupFrequency(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	families := map[int]bool{
		15:   true,
		30:   true,
		60:   true,
		360:  true,
		720:  true,
		1440: true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Backup Frequency can only be '15', '30', '60', '360', '720' or '1440'"))
	}

	return warnings, errors
}

func getRedisConnectionString(redisHostName string, sslPort int32, accessKey string, enableSslPort bool) string {
	return fmt.Sprintf("%s:%d,password=%s,ssl=%t,abortConnect=False", redisHostName, sslPort, accessKey, enableSslPort)
}
