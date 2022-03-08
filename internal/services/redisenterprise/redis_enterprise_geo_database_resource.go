package redisenterprise

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2022-01-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2022-01-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceRedisEnterpriseGeoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRedisEnterpriseGeoDatabaseCreateUpdate,
		Read:   resourceRedisEnterpriseGeoDatabaseRead,
		Update: resourceRedisEnterpriseGeoDatabaseCreateUpdate,
		Delete: resourceRedisEnterpriseGeoDatabaseDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := redisenterprise.ParseDatabaseID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "default",
				ValidateFunc: validate.RedisEnterpriseDatabaseName,
			},

			"cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
			},

			"client_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.ProtocolEncrypted),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.ProtocolEncrypted),
					string(redisenterprise.ProtocolPlaintext),
				}, false),
			},

			"clustering_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.ClusteringPolicyOSSCluster),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.ClusteringPolicyEnterpriseCluster),
					string(redisenterprise.ClusteringPolicyOSSCluster),
				}, false),
			},

			"eviction_policy": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(redisenterprise.EvictionPolicyVolatileLRU),
				ValidateFunc: validation.StringInSlice([]string{
					string(redisenterprise.EvictionPolicyAllKeysLFU),
					string(redisenterprise.EvictionPolicyAllKeysLRU),
					string(redisenterprise.EvictionPolicyAllKeysRandom),
					string(redisenterprise.EvictionPolicyVolatileLRU),
					string(redisenterprise.EvictionPolicyVolatileLFU),
					string(redisenterprise.EvictionPolicyVolatileTTL),
					string(redisenterprise.EvictionPolicyVolatileRandom),
					string(redisenterprise.EvictionPolicyNoEviction),
				}, false),
			},

			"linked_database_id": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MaxItems: 5,
				Set:      pluginsdk.HashString,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: databases.ValidateDatabaseID,
				},
			},

			"linked_database_group_nickname": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "geoGroup",
			},

			"force_unlink_database_id": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: databases.ValidateDatabaseID,
				},
			},

			"redi_search_module_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"redi_search_module_args": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},

			"redi_search_module_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"port": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      10000,
				ValidateFunc: validation.IntBetween(0, 65353),
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
		},
		// CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
		//	oldLinkedDbRaw, newLinkedDbRaw := d.GetChange("linked_database_id")
		//	linkedListHasChange := d.HasChange("linked_database_id")
		//	oldLinkedDb := oldLinkedDbRaw.(*pluginsdk.Set).List()
		//	newLinkedDb := newLinkedDbRaw.(*pluginsdk.Set).List()
		//	unlinkList := d.Get("force_unlink_database_id").([]interface{})
		//
		//	js, _ := json.Marshal(linkedListHasChange)
		//	log.Printf("DDDDDhas changes%s", js)
		//	js1, _ := json.Marshal(oldLinkedDb)
		//	log.Printf("DDDDDold list%s", js1)
		//	js2, _ := json.Marshal(newLinkedDb)
		//	log.Printf("DDDDDnew list%s", js2)
		//	js3, _ := json.Marshal(unlinkList)
		//	log.Printf("DDDDDunlink list%s", js3)
		//
		//	oldItemList := make(map[string]bool)
		//	for _, oldItem := range oldLinkedDb {
		//		oldItemList[oldItem.(string)] = true
		//	}
		//
		//	for _, newItem := range newLinkedDb {
		//		oldItemList[newItem.(string)] = false
		//	}
		//
		//	for _, unlinkItem := range unlinkList {
		//		//if !oldItemList[unlinkItem.(string)] && linkedListHasChange {
		//		//	return fmt.Errorf("The unlinked database must be a linked database and be removed from the linked database list")
		//		//}
		//		oldItemList[unlinkItem.(string)] = false
		//	}
		//
		//	for _, oldItem := range oldLinkedDb {
		//		if oldItemList[oldItem.(string)] {
		//			return fmt.Errorf("Please use forceUnlink action to remove a linked database from the list")
		//		}
		//	}
		//
		//	return nil
		// }),
	}
}

func resourceRedisEnterpriseGeoDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).RedisEnterprise.GeoDatabaseClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := redisenterprise.ParseRedisEnterpriseID(d.Get("cluster_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `cluster_id`: %+v", err)
	}

	id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redis_enterprise_geo_database", id.ID())
		}
	}

	clusteringPolicy := databases.ClusteringPolicy(d.Get("clustering_policy").(string))
	evictionPolicy := databases.EvictionPolicy(d.Get("eviction_policy").(string))
	protocol := databases.Protocol(d.Get("client_protocol").(string))

	// linkedDbHasChange := d.HasChange("linked_database_id")
	unlinkedDbList, hasUnlinkedDb := d.GetOk("force_unlink_database_id")
	if hasUnlinkedDb {
		oldLinkedDbRaw, _ := d.GetChange("linked_database_id")
		oldLinkedDb := oldLinkedDbRaw.(*pluginsdk.Set).List()
		if err := forceUnlinkDatabase(d, meta, oldLinkedDb, unlinkedDbList.([]interface{})); err != nil {
			return fmt.Errorf("unlinking database error: %+v", err)
		}
	}

	linkedDatabase, err := expandArmGeoLinkedDatabase(d.Get("linked_database_id").(*pluginsdk.Set).List(), id.ID(), d.Get("linked_database_group_nickname").(string))
	if err != nil {
		return fmt.Errorf("Setting geo database for database %s error: %+v", id.ID(), err)
	}

	parameters := databases.Database{
		Properties: &databases.DatabaseProperties{
			ClientProtocol:   &protocol,
			ClusteringPolicy: &clusteringPolicy,
			EvictionPolicy:   &evictionPolicy,
			GeoReplication:   linkedDatabase,
			Port:             utils.Int64(int64(d.Get("port").(int))),
		},
	}

	if d.Get("redi_search_module_enabled").(bool) {
		if evictionPolicy != databases.EvictionPolicyNoEviction {
			return fmt.Errorf("evictionPolicy must be set to NoEviction when using RediSearch module")
		}
		parameters.Properties.Modules = &[]databases.Module{
			{
				Name: "RediSearch",
				Args: utils.String(d.Get("redi_search_module_args").(string)),
			},
		}
	}

	future, err := client.Create(ctx, id, parameters)
	if err != nil {
		// @tombuildsstuff: investigate moving this above

		// Need to check if this was due to the cluster having the wrong sku
		// if strings.Contains(err.Error(), "The value of the parameter 'properties.modules' is invalid") {
		//	clusterClient := meta.(*clients.Client).RedisEnterprise.Client
		//	resp, err := clusterClient.Get(ctx, *clusterId)
		//	if err != nil {
		//		return fmt.Errorf("retrieving %s: %+v", *clusterId, err)
		//	}
		//
		//	if strings.Contains(strings.ToLower(string(resp.Model.Sku.Name)), "flash") {
		//		return fmt.Errorf("creating a Redis Enterprise Database with modules in a Redis Enterprise Cluster that has an incompatible Flash SKU type %q - please remove the Redis Enterprise Database modules or change the Redis Enterprise Cluster SKU type %s", string(resp.Model.Sku.Name), id)
		//	}
		// }

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRedisEnterpriseGeoDatabaseRead(d, meta)
}

func resourceRedisEnterpriseGeoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.GeoDatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Redis Enterprise Database %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keysResp, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("name", id.DatabaseName)
	clusterId := redisenterprise.NewRedisEnterpriseID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName)
	d.Set("cluster_id", clusterId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			clientProtocol := ""
			if props.ClientProtocol != nil {
				clientProtocol = string(*props.ClientProtocol)
			}
			d.Set("client_protocol", clientProtocol)

			clusteringPolicy := ""
			if props.ClusteringPolicy != nil {
				clusteringPolicy = string(*props.ClusteringPolicy)
			}
			d.Set("clustering_policy", clusteringPolicy)

			evictionPolicy := ""
			if props.EvictionPolicy != nil {
				evictionPolicy = string(*props.EvictionPolicy)
			}
			d.Set("eviction_policy", evictionPolicy)

			moduleEnabled := false
			var version string
			var args string
			if modules := props.Modules; modules != nil {
				for _, item := range *modules {
					if item.Args != nil {
						args = *item.Args
						d.Set("redi_search_module_args", args)
					}

					if item.Name == "RediSearch" {
						moduleEnabled = true
						d.Set("redi_search_module_enabled", moduleEnabled)
					}

					if item.Version != nil {
						version = *item.Version
						d.Set("redi_search_module_version", version)
					}
				}
			}
			d.Set("redi_search_module_enabled", moduleEnabled)
			d.Set("redi_search_module_args", args)
			d.Set("redi_search_module_version", version)

			if geoProps := props.GeoReplication; geoProps != nil {
				if geoProps.GroupNickname != nil {
					d.Set("linked_database_group_nickname", geoProps.GroupNickname)
				}
				if err := d.Set("linked_database_id", flattenArmGeoLinkedDatabase(geoProps.LinkedDatabases)); err != nil {
					return fmt.Errorf("setting `linked_database_id`: %+v", err)
				}
			}
			d.Set("port", props.Port)
		}
	}

	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	d.Set("force_unlink_database_id", d.Get("force_unlink_database_id").([]interface{}))

	return nil
}

func resourceRedisEnterpriseGeoDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.GeoDatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandArmGeoLinkedDatabase(inputId []interface{}, parentDBId string, inputGeoName string) (*databases.DatabasePropertiesGeoReplication, error) {
	idList := make([]databases.LinkedDatabase, 0)
	isParentDbIncluded := false

	for _, id := range inputId {
		if id.(string) == parentDBId {
			isParentDbIncluded = true
		}
		idList = append(idList, databases.LinkedDatabase{
			Id: utils.String(id.(string)),
		})
	}
	if isParentDbIncluded {
		return &databases.DatabasePropertiesGeoReplication{
			LinkedDatabases: &idList,
			GroupNickname:   utils.String(inputGeoName),
		}, nil
	}

	return nil, fmt.Errorf("linked database list must include database ID: %s", parentDBId)
}

func flattenArmGeoLinkedDatabase(inputDB *[]databases.LinkedDatabase) []string {
	results := make([]string, 0)

	if inputDB == nil {
		return results
	}

	for _, item := range *inputDB {
		if item.Id != nil {
			results = append(results, *item.Id)
		}
	}
	return results
}

func forceUnlinkDatabase(d *pluginsdk.ResourceData, meta interface{}, linkedDatabaseList []interface{}, unlinkedDbRaw []interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.GeoDatabaseClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO]Preparing to unlink a linked database")

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	linkedDb := make(map[string]bool)
	for _, linkedItem := range linkedDatabaseList {
		linkedDb[linkedItem.(string)] = true
	}

	for _, unlinkedItem := range unlinkedDbRaw {
		if !linkedDb[unlinkedItem.(string)] {
			return fmt.Errorf("%s is not a linked database", unlinkedItem)
		}
	}
	unlinkedDbList := utils.ExpandStringSlice(unlinkedDbRaw)
	parameters := databases.ForceUnlinkParameters{
		Ids: *unlinkedDbList,
	}

	if err := client.ForceUnlinkThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("force unlinking from database %s error: %+v", id, err)
	}

	return nil
}
