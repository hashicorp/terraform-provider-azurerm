// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedredis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// Geo-replication / linked databases are managed as a separate resource because when dbs are linked, ARM
// will mutate the state of all dbs out of bound, causing unexpected plan diff.
//
// The default database name is always "default", and because cluster and database are managed as a single TF resource,
// we use cluster id to configure linking. Internally the database id can be derived from cluster id
// by appending "/databases/default" suffix.

type ManagedRedisGeoReplicationResource struct{}

var (
	_ sdk.ResourceWithCustomizeDiff = ManagedRedisGeoReplicationResource{}
	_ sdk.ResourceWithUpdate        = ManagedRedisGeoReplicationResource{}
)

type ManagedRedisGeoReplicationResourceModel struct {
	ManagedRedisId        string   `tfschema:"managed_redis_id"`
	LinkedManagedRedisIds []string `tfschema:"linked_managed_redis_ids"`
}

func (r ManagedRedisGeoReplicationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_redis_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
		},

		"linked_managed_redis_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 4,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: redisenterprise.ValidateRedisEnterpriseID,
			},
		},
	}
}

func (r ManagedRedisGeoReplicationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedRedisGeoReplicationResource) ModelObject() interface{} {
	return &ManagedRedisGeoReplicationResourceModel{}
}

func (r ManagedRedisGeoReplicationResource) ResourceType() string {
	return "azurerm_managed_redis_geo_replication"
}

func (r ManagedRedisGeoReplicationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redisenterprise.ValidateRedisEnterpriseID
}

func (r ManagedRedisGeoReplicationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			var model ManagedRedisGeoReplicationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := redisenterprise.ParseRedisEnterpriseID(model.ManagedRedisId)
			if err != nil {
				return err
			}

			if err := linkUnlinkGeoReplication(ctx, client, model, id); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedRedisGeoReplicationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			resp, err := client.Get(ctx, dbId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(dbId)
				}
				return fmt.Errorf("retrieving %s: %+v", dbId, err)
			}

			state := ManagedRedisGeoReplicationResourceModel{
				ManagedRedisId: clusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil && props.GeoReplication != nil {
					state.LinkedManagedRedisIds = make([]string, 0, len(pointer.From(props.GeoReplication.LinkedDatabases)))
					for _, db := range pointer.From(props.GeoReplication.LinkedDatabases) {
						cId, err := toClusterId(pointer.From(db.Id))
						if err != nil {
							return err
						}
						if !resourceids.Match(cId, clusterId) {
							state.LinkedManagedRedisIds = append(state.LinkedManagedRedisIds, cId.ID())
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedRedisGeoReplicationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagedRedisGeoReplicationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := linkUnlinkGeoReplication(ctx, client, model, clusterId); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r ManagedRedisGeoReplicationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			existing, err := client.Get(ctx, dbId)
			if err != nil {
				return err
			}

			if existing.Model.Properties != nil && existing.Model.Properties.GeoReplication != nil {
				fromDbIds := flattenLinkedDatabases(existing.Model.Properties.GeoReplication.LinkedDatabases)
				toDbIds := []string{dbId.ID()}

				dbIdsToUnlink := databaselink.DbIdsToUnlink(fromDbIds, toDbIds)

				if len(dbIdsToUnlink) > 0 {
					params := databases.ForceUnlinkParameters{
						Ids: dbIdsToUnlink,
					}
					if err := client.ForceUnlinkThenPoll(ctx, dbId, params); err != nil {
						return fmt.Errorf("force unlink %s: %+v", dbId, err)
					}
				}
			}

			return nil
		},
	}
}

func (r ManagedRedisGeoReplicationResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			var model ManagedRedisGeoReplicationResourceModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return err
			}

			return nil
		},
	}
}

func linkUnlinkGeoReplication(ctx context.Context, client *databases.DatabasesClient, model ManagedRedisGeoReplicationResourceModel, clusterId *redisenterprise.RedisEnterpriseId) error {
	id := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		return err
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if existing.Model.Properties.GeoReplication == nil {
		return fmt.Errorf("geo_replication_group_name has to be set on database %s", id)
	}

	fromDbIds := flattenLinkedDatabases(existing.Model.Properties.GeoReplication.LinkedDatabases)
	toDbIds, err := toDbIds(model.LinkedManagedRedisIds, id)
	if err != nil {
		return err
	}

	dbIdsToUnlink := databaselink.DbIdsToUnlink(fromDbIds, toDbIds)

	if len(dbIdsToUnlink) > 0 {
		params := databases.ForceUnlinkParameters{
			Ids: dbIdsToUnlink,
		}
		if err := client.ForceUnlinkThenPoll(ctx, id, params); err != nil {
			return fmt.Errorf("force unlink %s: %+v", id, err)
		}
	}

	if databaselink.HasDbToLink(fromDbIds, toDbIds) {
		params := databases.ForceLinkParameters{
			GeoReplication: databases.ForceLinkParametersGeoReplication{
				GroupNickname:   existing.Model.Properties.GeoReplication.GroupNickname,
				LinkedDatabases: expandLinkedDatabases(toDbIds),
			},
		}

		err = client.ForceLinkToReplicationGroupThenPoll(ctx, id, params)
		if err != nil {
			return fmt.Errorf("force link %s: %+v", id, err)
		}
	}
	return nil
}

func toDbIds(otherClusterIds []string, selfDbId databases.DatabaseId) ([]string, error) {
	dbIds := make([]string, 0, len(otherClusterIds)+1)
	containsSelf := false

	for _, cIdStr := range otherClusterIds {
		cId, err := redisenterprise.ParseRedisEnterpriseID(cIdStr)
		if err != nil {
			return nil, err
		}
		otherDbId := databases.NewDatabaseID(cId.SubscriptionId, cId.ResourceGroupName, cId.RedisEnterpriseName, defaultDatabaseName)

		if resourceids.Match(&otherDbId, &selfDbId) {
			containsSelf = true
		}

		dbIds = append(dbIds, otherDbId.ID())
	}

	if !containsSelf {
		dbIds = append(dbIds, selfDbId.ID())
	}

	return dbIds, nil
}

func toClusterId(dbIdStr string) (*redisenterprise.RedisEnterpriseId, error) {
	dbId, err := databases.ParseDatabaseID(dbIdStr)
	if err != nil {
		return nil, err
	}
	return pointer.To(redisenterprise.NewRedisEnterpriseID(dbId.SubscriptionId, dbId.ResourceGroupName, dbId.RedisEnterpriseName)), nil
}

func flattenLinkedDatabases(dbs *[]databases.LinkedDatabase) []string {
	if dbs == nil {
		return nil
	}

	result := make([]string, 0, len(*dbs))
	for _, db := range *dbs {
		if db.Id != nil {
			result = append(result, pointer.From(db.Id))
		}
	}
	return result
}

func expandLinkedDatabases(dbIds []string) *[]databases.LinkedDatabase {
	if len(dbIds) == 0 {
		return nil
	}

	result := make([]databases.LinkedDatabase, 0, len(dbIds))
	for _, id := range dbIds {
		result = append(result, databases.LinkedDatabase{
			Id: pointer.To(id),
		})
	}
	return &result
}
