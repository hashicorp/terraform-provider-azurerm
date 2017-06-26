package azurerm

import (
	"fmt"
	"log"

	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/redis"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/jen20/riviera/azure"
)

func resourceArmRedisCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRedisCacheCreate,
		Read:   resourceArmRedisCacheRead,
		Update: resourceArmRedisCacheUpdate,
		Delete: resourceArmRedisCacheDelete,
		// TODO: import support
		MigrateState:  resourceAzureRMRedisCacheMigrateState,
		SchemaVersion: 1,

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

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateRedisSku,
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

			"redis_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maxclients": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"maxmemory_delta": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"maxmemory_reserved": {
							Type:     schema.TypeString,
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
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "false",
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"true",
								"false",
							}, true),
						},
						"rdb_backup_frequency": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "",
							ValidateFunc: validateRedisBackupFrequency,
						},
						"rdb_backup_max_snapshot_count": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"rdb_storage_connection_string": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
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
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmRedisCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient
	log.Printf("[INFO] preparing arguments for Azure ARM Redis Cache creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)

	enableNonSSLPort := d.Get("enable_non_ssl_port").(bool)

	capacity := int32(d.Get("capacity").(int))
	family := redis.SkuFamily(d.Get("family").(string))
	sku := redis.SkuName(d.Get("sku_name").(string))

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := redis.CreateParameters{
		Name:     &name,
		Location: &location,
		CreateProperties: &redis.CreateProperties{
			EnableNonSslPort: &enableNonSSLPort,
			Sku: &redis.Sku{
				Capacity: &capacity,
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

	_, error := client.Create(resGroup, name, parameters, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
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
		Refresh:    redisStateRefreshFunc(client, resGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Instance (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	return resourceArmRedisCacheRead(d, meta)
}

func resourceArmRedisCacheUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient
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
			EnableNonSslPort: &enableNonSSLPort,
			Sku: &redis.Sku{
				Capacity: &capacity,
				Family:   family,
				Name:     sku,
			},
			Tags: expandedTags,
		},
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

	_, err := client.Update(resGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
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
		Refresh:    redisStateRefreshFunc(client, resGroup, name),
		Timeout:    60 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Redis Instance (%s) to become available: %s", d.Get("name"), err)
	}

	d.SetId(*read.ID)

	return resourceArmRedisCacheRead(d, meta)
}

func resourceArmRedisCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	resp, err := client.Get(resGroup, name)

	// covers if the resource has been deleted outside of TF, but is still in the state
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Redis Cache %s: %s", name, err)
	}

	keysResp, err := client.ListKeys(resGroup, name)
	if err != nil {
		return fmt.Errorf("Error making ListKeys request on Azure Redis Cache %s: %s", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("ssl_port", resp.SslPort)
	d.Set("hostname", resp.HostName)
	d.Set("port", resp.Port)
	d.Set("enable_non_ssl_port", resp.EnableNonSslPort)
	d.Set("capacity", resp.Sku.Capacity)
	d.Set("family", resp.Sku.Family)
	d.Set("sku_name", resp.Sku.Name)

	if resp.ShardCount != nil {
		d.Set("shard_count", resp.ShardCount)
	}

	redisConfiguration := flattenRedisConfiguration(resp.RedisConfiguration)
	d.Set("redis_configuration", &redisConfiguration)

	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRedisCacheDelete(d *schema.ResourceData, meta interface{}) error {
	redisClient := meta.(*ArmClient).redisClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["Redis"]

	deleteResp, error := redisClient.Delete(resGroup, name, make(chan struct{}))
	resp := <-deleteResp
	err = <-error

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error issuing Azure ARM delete request of Redis Cache Instance '%s': %s", name, err)
	}

	checkResp, _ := redisClient.Get(resGroup, name)
	if checkResp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("Error issuing Azure ARM delete request of Redis Cache Instance '%s': it still exists after deletion", name)
	}

	return nil
}

func redisStateRefreshFunc(client redis.GroupClient, resourceGroupName string, sgName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(resourceGroupName, sgName)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in redisStateRefreshFunc to Azure ARM for Redis Cache Instance '%s' (RG: '%s'): %s", sgName, resourceGroupName, err)
		}

		return res, *res.ProvisioningState, nil
	}
}

func expandRedisConfiguration(d *schema.ResourceData) *map[string]*string {
	configuration := d.Get("redis_configuration").([]interface{})

	output := make(map[string]*string)

	if configuration == nil {
		return &output
	}

	// TODO: can we use this to remove the below? \/
	//config := configuration[0].(map[string]interface{})

	for _, v := range configuration {
		config := v.(map[string]interface{})

		maxClients := config["maxclients"].(string)
		if maxClients != "" {
			output["maxclients"] = azure.String(maxClients)
		}

		maxMemoryDelta := config["maxmemory_delta"].(string)
		if maxMemoryDelta != "" {
			output["maxmemory-delta"] = azure.String(maxMemoryDelta)
		}

		maxMemoryReserved := config["maxmemory_reserved"].(string)
		if maxMemoryReserved != "" {
			output["maxmemory-reserved"] = azure.String(maxMemoryReserved)
		}

		maxMemoryPolicy := config["maxmemory_policy"].(string)
		if maxMemoryPolicy != "" {
			output["maxmemory-policy"] = azure.String(maxMemoryPolicy)
		}

		backupEnabled := config["rdb_backup_enabled"].(string)
		output["rdb-backup-enabled"] = azure.String(backupEnabled)

		backupFrequency := config["rdb_backup_frequency"].(string)
		if backupFrequency != "" {
			output["rdb-backup-frequency"] = azure.String(backupFrequency)
		}

		backupMaxSnapshotCount := config["rdb_backup_max_snapshot_count"].(string)
		if backupMaxSnapshotCount != "" {
			output["rdb-backup-max-snapshot-count"] = azure.String(backupMaxSnapshotCount)
		}

		backupStorageConnectionString := config["rdb_storage_connection_string"].(string)
		if backupStorageConnectionString != "" {
			output["rdb-storage-connection-string"] = azure.String(backupStorageConnectionString)
		}
	}

	return &output
}

func flattenRedisConfiguration(configuration *map[string]*string) map[string]*string {
	redisConfiguration := make(map[string]*string, len(*configuration))
	config := *configuration

	redisConfiguration["maxclients"] = config["maxclients"]
	redisConfiguration["maxmemory_delta"] = config["maxmemory-delta"]
	redisConfiguration["maxmemory_reserved"] = config["maxmemory-reserved"]
	redisConfiguration["maxmemory_policy"] = config["maxmemory-policy"]

	redisConfiguration["rdb_backup_enabled"] = config["rdb-backup-enabled"]
	redisConfiguration["rdb_backup_frequency"] = config["rdb-backup-frequency"]
	redisConfiguration["rdb_backup_max_snapshot_count"] = config["rdb-backup-max-snapshot-count"]
	redisConfiguration["rdb_storage_connection_string"] = config["rdb-storage-connection-string"]

	return redisConfiguration
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

func validateRedisSku(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	skus := map[string]bool{
		"basic":    true,
		"standard": true,
		"premium":  true,
	}

	if !skus[value] {
		errors = append(errors, fmt.Errorf("Redis SKU can only be Basic, Standard or Premium"))
	}
	return
}

func validateRedisBackupFrequency(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	families := map[string]bool{
		"15":   true,
		"30":   true,
		"60":   true,
		"360":  true,
		"720":  true,
		"1440": true,
	}

	if !families[value] {
		errors = append(errors, fmt.Errorf("Redis Backup Frequency can only be '15', '30', '60', '360', '720' or '1440'"))
	}

	return
}
