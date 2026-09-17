// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_replication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func (r Resource) Delete() sdk.ResourceFunc {
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

				dbIdsToUnlink, intermediateDbIds, _ := databaselink.LinkUnlink(fromDbIds, toDbIds)

				for _, inv := range databaselink.ForceUnlinkInvocations(intermediateDbIds, dbIdsToUnlink) {
					id, err := databases.ParseDatabaseID(inv.Id)
					if err != nil {
						return err
					}

					params := databases.ForceUnlinkParameters{
						Ids: inv.Ids,
					}

					if err := client.ForceUnlinkThenPoll(ctx, *id, params); err != nil {
						return fmt.Errorf("force unlink %s: %+v", *id, err)
					}
				}
			}

			return nil
		},
	}
}

func flattenLinkedDatabases(dbs *[]databases.LinkedDatabase) []string {
	if dbs == nil {
		return []string{}
	}

	result := make([]string, 0, len(*dbs))
	for _, db := range *dbs {
		if db.Id != nil {
			result = append(result, pointer.From(db.Id))
		}
	}
	return result
}
