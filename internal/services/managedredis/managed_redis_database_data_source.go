// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedRedisDatabaseDataSource struct{}

var _ sdk.DataSource = ManagedRedisDatabaseDataSource{}

type ManagedRedisDatabaseDataSourceModel struct {
	Name                        string   `tfschema:"name"`
	ClusterId                   string   `tfschema:"cluster_id"`
	LinkedDatabaseId            []string `tfschema:"linked_database_id"`
	LinkedDatabaseGroupNickname string   `tfschema:"linked_database_group_nickname"`
	PrimaryAccessKey            string   `tfschema:"primary_access_key"`
	SecondaryAccessKey          string   `tfschema:"secondary_access_key"`
}

func (r ManagedRedisDatabaseDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
		},
	}
}

func (r ManagedRedisDatabaseDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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
}

func (r ManagedRedisDatabaseDataSource) ModelObject() interface{} {
	return &ManagedRedisDatabaseDataSourceModel{}
}

func (r ManagedRedisDatabaseDataSource) ResourceType() string {
	return "azurerm_managed_redis_database"
}

func (r ManagedRedisDatabaseDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagedRedisDatabaseDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(model.ClusterId)
			if err != nil {
				return err
			}

			id := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, model.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ManagedRedisDatabaseDataSourceModel{
				Name:      id.DatabaseName,
				ClusterId: clusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.GeoReplication != nil {
						if props.GeoReplication.GroupNickname != nil {
							state.LinkedDatabaseGroupNickname = *props.GeoReplication.GroupNickname
						}
						if props.GeoReplication.LinkedDatabases != nil {
							state.LinkedDatabaseId = flattenArmGeoLinkedDatabase(props.GeoReplication.LinkedDatabases)
						}
					}
					if strings.EqualFold(string(pointer.From(props.AccessKeysAuthentication)), string(databases.AccessKeysAuthenticationEnabled)) {
						keysResp, err := client.ListKeys(ctx, id)
						if err != nil {
							return fmt.Errorf("listing keys for %s: %+v", id, err)
						}
						if keysModel := keysResp.Model; keysModel != nil {
							state.PrimaryAccessKey = pointer.From(keysModel.PrimaryKey)
							state.SecondaryAccessKey = pointer.From(keysModel.SecondaryKey)
						}
					}
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
