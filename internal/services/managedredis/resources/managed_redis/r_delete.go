// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient

			clusterId, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			dbId := databases.NewDatabaseID(clusterId.SubscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			if err := dbClient.DeleteThenPoll(ctx, dbId); err != nil {
				return fmt.Errorf("deleting %s: %+v", dbId, err)
			}

			if err := clusterClient.DeleteThenPoll(ctx, *clusterId); err != nil {
				return fmt.Errorf("deleting %s: %+v", clusterId, err)
			}

			return nil
		},
	}
}
