// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redisenterprise

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2024-10-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2024-10-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceRedisEnterpriseDatabase() *pluginsdk.Resource {
	s := map[string]*pluginsdk.Schema{
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
	}

	return &pluginsdk.Resource{
		Read: dataSourceRedisEnterpriseDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: s,
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

	id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, d.Get("name").(string))

	keysResp, err := client.ListKeys(ctx, id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", id, err)
	}

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.DatabaseName)
	d.Set("cluster_id", clusterId.ID())

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

	if model := keysResp.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("secondary_access_key", model.SecondaryKey)
	}

	return nil
}
