package redisenterprise

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2022-01-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2022-01-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceRedisEnterpriseGeoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceRedisEnterpriseGepDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
			},

			"linked_database_id": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: databases.ValidateDatabaseID,
				},
			},

			"linked_database_group_nickname": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "geoGroup",
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
	}
}

func dataSourceRedisEnterpriseGepDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.GeoDatabaseClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := redisenterprise.ParseRedisEnterpriseID(d.Get("cluster_id").(string))
	if err != nil {
		return err
	}

	id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.GeoReplication != nil {
			if props.GeoReplication.GroupNickname != nil {
				d.Set("linked_database_group_nickname", props.GeoReplication.GroupNickname)
			}
			if props.GeoReplication.LinkedDatabases != nil {
				d.Set("linked_database_id", flattenArmGeoLinkedDatabase(props.GeoReplication.LinkedDatabases))
			}
		}
	}

	keysResp, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}
	d.SetId(d.Id())

	d.Set("name", id.DatabaseName)
	d.Set("cluster_id", clusterId.ID())

	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}
