// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis_geo_replication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
						if pointer.From(db.State) == databases.LinkStateLinked {
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
			}

			return metadata.Encode(&state)
		},
	}
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
