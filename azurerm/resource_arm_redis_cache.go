package azurerm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
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
				StateFunc: azureRMNormalizeLocation,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"family": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateRedisFamily,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.Basic),
					string(redis.Standard),
					string(redis.Premium),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
				Required: true,
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"Monday",
								"Tuesday",
								"Wednesday",
								"Thursday",
								"Friday",
								"Saturday",
								"Sunday",
							}, true),
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmRedisCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	patchSchedule, err := expandRedisPatchSchedule(d)
	if err != nil {
		return fmt.Errorf("Error parsing Patch Schedule: %+v", err)
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
			RedisConfiguration: expandRedisConfiguration(d),
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
		parameters.SubnetID = utils.String(v.(string))
	}

	future, err := client.Create(ctx, resGroup, name, parameters)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
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
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, resGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Instance (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	if schedule := patchSchedule; schedule != nil {
		patchClient := meta.(*ArmClient).redisPatchSchedulesClient
		_, err = patchClient.CreateOrUpdate(ctx, resGroup, name, *schedule)
		if err != nil {
			return fmt.Errorf("Error setting Redis Patch Schedule: %+v", err)
		}
	}

	return resourceArmRedisCacheRead(d, meta)
}

func resourceArmRedisCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache update.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := redis.UpdateParameters{
		UpdateProperties: &redis.UpdateProperties{
			EnableNonSslPort: utils.Bool(enableNonSSLPort),
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
		redisConfiguration := expandRedisConfiguration(d)
		parameters.RedisConfiguration = redisConfiguration
	}

	_, err := client.Update(ctx, resGroup, name, parameters)
	if err != nil {
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
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, resGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Instance (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	patchSchedule, err := expandRedisPatchSchedule(d)
	if err != nil {
		return fmt.Errorf("Error parsing Patch Schedule: %+v", err)
	}

	patchClient := meta.(*ArmClient).redisPatchSchedulesClient
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
	client := meta.(*ArmClient).redisClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	patchSchedulesClient := meta.(*ArmClient).redisPatchSchedulesClient

	schedule, err := patchSchedulesClient.Get(ctx, resGroup, name)
	if err == nil {
		patchSchedule := flattenRedisPatchSchedules(schedule)
		if err := d.Set("patch_schedule", patchSchedule); err != nil {
			return fmt.Errorf("Error setting `patch_schedule`: %+v", err)
		}
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("capacity", sku.Capacity)
		d.Set("family", sku.Family)
		d.Set("sku_name", sku.Name)
	}

	if props := resp.Properties; props != nil {
		d.Set("ssl_port", props.SslPort)
		d.Set("hostname", props.HostName)
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

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRedisCacheDelete(d *schema.ResourceData, meta interface{}) error {
	redisClient := meta.(*ArmClient).redisClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	future, err := redisClient.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}
	err = future.WaitForCompletion(ctx, redisClient.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	return nil
}

func redisStateRefreshFunc(ctx context.Context, client redis.Client, resourceGroupName string, sgName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, sgName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in redisStateRefreshFunc to Azure ARM for Redis Cache Instance '%s' (RG: '%s'): %s", sgName, resourceGroupName, err)
		}

		return res, string(res.ProvisioningState), nil
	}
}

func expandRedisConfiguration(d *schema.ResourceData) map[string]*string {
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

	// Backup
	if v, ok := d.GetOk("redis_configuration.0.rdb_backup_enabled"); ok {
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

	return output
}

func expandRedisPatchSchedule(d *schema.ResourceData) (*redis.PatchSchedule, error) {
	v, ok := d.GetOk("patch_schedule")
	if !ok {
		return nil, nil
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
	return &schedule, nil
}

func flattenRedisConfiguration(input map[string]*string) ([]interface{}, error) {
	outputs := make(map[string]interface{}, len(input))

	if v := input["maxclients"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxclients` %q: %+v", v, err)
		}
		outputs["maxclients"] = i
	}
	if v := input["maxmemory-delta"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxmemory-delta` %q: %+v", v, err)
		}
		outputs["maxmemory_delta"] = i
	}
	if v := input["maxmemory-reserved"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `maxmemory-reserved` %q: %+v", v, err)
		}
		outputs["maxmemory_reserved"] = i
	}
	if v := input["maxmemory-policy"]; v != nil {
		outputs["maxmemory_policy"] = *v
	}

	// delta, reserved, enabled, frequency,, count,
	if v := input["rdb-backup-enabled"]; v != nil {
		b, err := strconv.ParseBool(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-enabled` %q: %+v", v, err)
		}
		outputs["rdb_backup_enabled"] = b
	}
	if v := input["rdb-backup-frequency"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-frequency` %q: %+v", v, err)
		}
		outputs["rdb_backup_frequency"] = i
	}
	if v := input["rdb-backup-max-snapshot-count"]; v != nil {
		i, err := strconv.Atoi(*v)
		if err != nil {
			return nil, fmt.Errorf("Error parsing `rdb-backup-max-snapshot-count` %q: %+v", v, err)
		}
		outputs["rdb_backup_max_snapshot_count"] = i
	}
	if v := input["rdb-storage-connection-string"]; v != nil {
		outputs["rdb_storage_connection_string"] = *v
	}
	if v := input["notify-keyspace-events"]; v != nil {
		outputs["notify_keyspace_events"] = *v
	}

	return []interface{}{outputs}, nil
}

func flattenRedisPatchSchedules(schedule redis.PatchSchedule) []interface{} {
	outputs := make([]interface{}, 0)

	for _, entry := range *schedule.ScheduleEntries.ScheduleEntries {
		output := make(map[string]interface{}, 0)

		output["day_of_week"] = string(entry.DayOfWeek)
		output["start_hour_utc"] = int(*entry.StartHourUtc)

		outputs = append(outputs, output)
	}

	return outputs
}

func validateRedisFamily(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"c": true,
		"p": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Family can only be C or P"))
	}
	return
}

func validateRedisMaxMemoryPolicy(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	families := map[string]bool{
		"noeviction":      true,
		"allkeys-lru":     true,
		"volatile-lru":    true,
		"allkeys-random":  true,
		"volatile-random": true,
		"volatile-ttl":    true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Max Memory Policy can only be 'noeviction' / 'allkeys-lru' / 'volatile-lru' / 'allkeys-random' / 'volatile-random' / 'volatile-ttl'"))
	}

	return
}

func validateRedisBackupFrequency(v interface{}, k string) (ws []string, errors []error) {
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

	return
}
