// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_replication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/databaselink"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := linkUnlinkGeoReplication(ctx, metadata, client, model, clusterId, false); err != nil {
				return err
			}

			return nil
		},
	}
}

func linkUnlinkGeoReplication(ctx context.Context, metadata sdk.ResourceMetaData, client *databases.DatabasesClient, model ResourceModel, clusterId *redisenterprise.RedisEnterpriseId, create bool) error {
	primaryId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

	existing, err := client.Get(ctx, primaryId)
	if err != nil {
		return err
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", primaryId)
	}
	if existing.Model.Properties.GeoReplication == nil {
		return fmt.Errorf("geo_replication_group_name has to be set on database %s", primaryId)
	}

	fromDbIds := flattenLinkedDatabases(existing.Model.Properties.GeoReplication.LinkedDatabases)
	toDbIds, err := toDbIds(model.LinkedManagedRedisIds, primaryId)
	if err != nil {
		return err
	}

	dbIdsToUnlink, intermediateDbIds, dbIdsToLink := databaselink.LinkUnlink(fromDbIds, toDbIds)

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

		// Workaround for race-condition bug after force-unlinking
		// The API bug will be fixed in https://github.com/Azure/azure-rest-api-specs/issues/39598
		pollerType := custompollers.NewGeoReplicationUnlinkingPoller(client, primaryId, inv.Ids)
		poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for `linked_managed_redis_id` state to be consistent after unlinking for %s: %+v", primaryId, err)
		}
	}

	for _, inv := range databaselink.ForceLinkInvocations(intermediateDbIds, dbIdsToLink) {
		id, err := databases.ParseDatabaseID(inv.Id)
		if err != nil {
			return err
		}

		params := databases.ForceLinkParameters{
			GeoReplication: databases.ForceLinkParametersGeoReplication{
				GroupNickname:   existing.Model.Properties.GeoReplication.GroupNickname,
				LinkedDatabases: expandLinkedDatabases(inv.Ids),
			},
		}

		if create {
			if err := client.ForceLinkToReplicationGroupCallbackThenPoll(ctx, *id, params, metadata.SetIDCallback(id)); err != nil {
				return fmt.Errorf("force link %s: %+v", *id, err)
			}
		} else {
			if err := client.ForceLinkToReplicationGroupThenPoll(ctx, *id, params); err != nil {
				return fmt.Errorf("force link %s: %+v", *id, err)
			}
		}

		// Workaround for race-condition bug after force-linking
		// The API bug will be fixed in https://github.com/Azure/azure-rest-api-specs/issues/39598
		pollerType := custompollers.NewGeoReplicationLinkingPoller(client, primaryId, inv.Ids)
		poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for `linked_managed_redis_id` state to be consistent after linking for %s: %+v", primaryId, err)
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
