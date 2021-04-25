package redisenterprise

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redisenterprise/mgmt/2021-03-01/redisenterprise"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redisenterprise/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redisenterprise/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceRedisEnterpriseDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisEnterpriseDatabaseCreate,
		Read:   resourceRedisEnterpriseDatabaseRead,
		// Update currently is not implemented, will be for GA
		Delete: resourceRedisEnterpriseDatabaseDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RedisEnterpriseDatabaseID(id)
			return err
		}),

		// Since update is not currently supported all attribute have to be marked as FORCE NEW
		// until support for Update comes online in the near future
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "default",
				ValidateFunc: validate.RedisEnterpriseDatabaseName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RedisEnterpriseClusterID,
			},

			"client_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.Encrypted),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.Encrypted),
					string(redisenterprise.Plaintext),
				}, false),
			},

			"clustering_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.OSSCluster),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.EnterpriseCluster),
					string(redisenterprise.OSSCluster),
				}, false),
			},

			"eviction_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.VolatileLRU),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.AllKeysLFU),
					string(redisenterprise.AllKeysLRU),
					string(redisenterprise.AllKeysRandom),
					string(redisenterprise.VolatileLRU),
					string(redisenterprise.VolatileLFU),
					string(redisenterprise.VolatileTTL),
					string(redisenterprise.VolatileRandom),
					string(redisenterprise.NoEviction),
				}, false),
			},

			"module": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"RedisBloom",
								"RedisTimeSeries",
								"RediSearch",
							}, false),
						},

						"args": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "",
						},

						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			// This attribute is currently in preview and is not returned by the RP
			// "persistence": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	MaxItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"aof_enabled": {
			// 				Type:     schema.TypeBool,
			// 				Optional: true,
			// 			},

			// 			"aof_frequency": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				ValidateFunc: validation.StringInSlice([]string{
			// 					string(redisenterprise.Ones),
			// 					string(redisenterprise.Always),
			// 				}, false),
			// 			},

			// 			"rdb_enabled": {
			// 				Type:     schema.TypeBool,
			// 				Optional: true,
			// 			},

			// 			"rdb_frequency": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 				ValidateFunc: validation.StringInSlice([]string{
			// 					string(redisenterprise.Oneh),
			// 					string(redisenterprise.Sixh),
			// 					string(redisenterprise.OneTwoh),
			// 				}, false),
			// 			},
			// 		},
			// 	},
			// },

			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      10000,
				ValidateFunc: validation.IntBetween(0, 65353),
			},
		},
	}
}
func resourceRedisEnterpriseDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).RedisEnterprise.DatabaseClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	clusterID, _ := parse.RedisEnterpriseClusterID(d.Get("cluster_id").(string))
	id := parse.NewRedisEnterpriseDatabaseID(subscriptionId, resourceGroup, clusterID.RedisEnterpriseName, name)

	existing, err := client.Get(ctx, resourceGroup, id.RedisEnterpriseName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", name, resourceGroup, id.RedisEnterpriseName, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_redis_enterprise_database", id.ID())
	}

	parameters := redisenterprise.Database{
		DatabaseProperties: &redisenterprise.DatabaseProperties{
			ClientProtocol:   redisenterprise.Protocol(d.Get("client_protocol").(string)),
			ClusteringPolicy: redisenterprise.ClusteringPolicy(d.Get("clustering_policy").(string)),
			EvictionPolicy:   redisenterprise.EvictionPolicy(d.Get("eviction_policy").(string)),
			Modules:          expandArmDatabaseModuleArray(d.Get("module").([]interface{})),
			//Persistence:      expandArmDatabasePersistence(d.Get("persistence").([]interface{})),
			Port: utils.Int32(int32(d.Get("port").(int))),
		},
	}

	future, err := client.Create(ctx, resourceGroup, id.RedisEnterpriseName, name, parameters)
	if err != nil {
		// Need to check if this was due to the cluster having the wrong sku
		if strings.Contains(err.Error(), "The value of the parameter 'properties.modules' is invalid") {
			clusterClient := meta.(*clients.Client).RedisEnterprise.Client
			resp, err := clusterClient.Get(ctx, clusterID.ResourceGroup, clusterID.RedisEnterpriseName)
			if err != nil {
				return fmt.Errorf("retrieving Redis Enterprise Cluster (Name %q / Resource Group %q): %+v", clusterID.RedisEnterpriseName, clusterID.ResourceGroup, err)
			}

			if strings.Contains(strings.ToLower(string(resp.Sku.Name)), "flash") {
				return fmt.Errorf("creating a Redis Enterprise Database with modules in a Redis Enterprise Cluster that has an incompatible Flash SKU type %q - please remove the Redis Enterprise Database modules or change the Redis Enterprise Cluster SKU type (Resource Group %q / Cluster Name %q / Database %q)", string(resp.Sku.Name), resourceGroup, id.RedisEnterpriseName, name)
			}
		}

		return fmt.Errorf("creating Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", name, resourceGroup, id.RedisEnterpriseName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creating future for Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", name, resourceGroup, id.RedisEnterpriseName, err)
	}

	d.SetId(id.ID())

	return resourceRedisEnterpriseDatabaseRead(d, meta)
}

func resourceRedisEnterpriseDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.DatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RedisEnterpriseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RedisEnterpriseName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Redis Enterprise Database %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", id.DatabaseName, id.ResourceGroup, id.RedisEnterpriseName, err)
	}

	d.Set("name", id.DatabaseName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_id", parse.NewRedisEnterpriseClusterID(id.SubscriptionId, id.ResourceGroup, id.RedisEnterpriseName).ID())

	if props := resp.DatabaseProperties; props != nil {
		d.Set("client_protocol", props.ClientProtocol)
		d.Set("clustering_policy", props.ClusteringPolicy)
		d.Set("eviction_policy", props.EvictionPolicy)
		if err := d.Set("module", flattenArmDatabaseModuleArray(props.Modules)); err != nil {
			return fmt.Errorf("setting `module`: %+v", err)
		}
		// if err := d.Set("persistence", flattenArmDatabasePersistence(props.Persistence)); err != nil {
		// 	return fmt.Errorf("setting `persistence`: %+v", err)
		// }
		d.Set("port", props.Port)
	}

	return nil
}

func resourceRedisEnterpriseDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.DatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.RedisEnterpriseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.RedisEnterpriseName, id.DatabaseName)
	if err != nil {
		return fmt.Errorf("deleting Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", id.DatabaseName, id.ResourceGroup, id.RedisEnterpriseName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting future for Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", id.DatabaseName, id.ResourceGroup, id.RedisEnterpriseName, err)
	}
	return nil
}

func expandArmDatabaseModuleArray(input []interface{}) *[]redisenterprise.Module {
	results := make([]redisenterprise.Module, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, redisenterprise.Module{
			Name: utils.String(v["name"].(string)),
			Args: utils.String(v["args"].(string)),
		})
	}
	return &results
}

// Persistence is currently preview and does not return from the RP but will be fully supported in the near future
// func expandArmDatabasePersistence(input []interface{}) *redisenterprise.Persistence {
// 	if len(input) == 0 {
// 		return nil
// 	}
// 	v := input[0].(map[string]interface{})
// 	return &redisenterprise.Persistence{
// 		AofEnabled:   utils.Bool(v["aof_enabled"].(bool)),
// 		AofFrequency: redisenterprise.AofFrequency(v["aof_frequency"].(string)),
// 		RdbEnabled:   utils.Bool(v["rdb_enabled"].(bool)),
// 		RdbFrequency: redisenterprise.RdbFrequency(v["rdb_frequency"].(string)),
// 	}
// }

func flattenArmDatabaseModuleArray(input *[]redisenterprise.Module) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		args := ""
		if item.Args != nil {
			args = *item.Args
		}

		var version string
		if item.Version != nil {
			version = *item.Version
		}

		results = append(results, map[string]interface{}{
			"name":    name,
			"args":    args,
			"version": version,
		})
	}

	return results
}

// Persistence is currently preview and does not return from the RP but will be fully supported in the near future
// func flattenArmDatabasePersistence(input *redisenterprise.Persistence) []interface{} {
// 	if input == nil {
// 		return make([]interface{}, 0)
// 	}

// 	var aofEnabled bool
// 	if input.AofEnabled != nil {
// 		aofEnabled = *input.AofEnabled
// 	}

// 	var aofFrequency redisenterprise.AofFrequency
// 	if input.AofFrequency != "" {
// 		aofFrequency = input.AofFrequency
// 	}

// 	var rdbEnabled bool
// 	if input.RdbEnabled != nil {
// 		rdbEnabled = *input.RdbEnabled
// 	}

// 	var rdbFrequency redisenterprise.RdbFrequency
// 	if input.RdbFrequency != "" {
// 		rdbFrequency = input.RdbFrequency
// 	}

// 	return []interface{}{
// 		map[string]interface{}{
// 			"aof_enabled":   aofEnabled,
// 			"aof_frequency": aofFrequency,
// 			"rdb_enabled":   rdbEnabled,
// 			"rdb_frequency": rdbFrequency,
// 		},
// 	}
// }
