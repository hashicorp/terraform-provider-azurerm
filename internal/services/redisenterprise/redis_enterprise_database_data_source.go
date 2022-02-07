package redisenterprise

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2021-08-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2021-08-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceRedisEnterpriseDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceRedisEnterpriseDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			// TODO: deprecate me
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
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

func dataSourceRedisEnterpriseDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RedisEnterprise.DatabaseClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	clusterId, err := redisenterprise.ParseRedisEnterpriseID(d.Get("cluster_id").(string))
	if err != nil {
		return err
	}

	id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.ClusterName, d.Get("name").(string))

	keysResp, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DatabaseName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_id", clusterId.ID())

	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}
