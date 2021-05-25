package redisenterprise

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redisenterprise/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redisenterprise/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"cluster_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RedisEnterpriseClusterID,
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

	clusterId, err := parse.RedisEnterpriseClusterID(d.Get("cluster_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewRedisEnterpriseDatabaseID(subscriptionId, d.Get("resource_group_name").(string), clusterId.RedisEnterpriseName, d.Get("name").(string))

	keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.RedisEnterpriseName, id.DatabaseName)
	if err != nil {
		return fmt.Errorf("listing keys for Redis Enterprise Database %q (Resource Group %q / Cluster Name %q): %+v", id.DatabaseName, id.ResourceGroup, id.RedisEnterpriseName, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.DatabaseName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_id", clusterId.ID())
	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	return nil
}
