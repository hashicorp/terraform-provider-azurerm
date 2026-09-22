// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-07-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedredis/custompollers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 45 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			clusterClient := metadata.Client.ManagedRedis.Client
			dbClient := metadata.Client.ManagedRedis.DatabaseClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagedRedisResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			clusterId := redisenterprise.NewRedisEnterpriseID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existingCluster, err := clusterClient.Get(ctx, clusterId)
				if err != nil {
					if !response.WasNotFound(existingCluster.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", clusterId, err)
					}
				}

				if !response.WasNotFound(existingCluster.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), clusterId)
				}
			}

			dbId := databases.NewDatabaseID(subscriptionId, clusterId.ResourceGroupName, clusterId.RedisEnterpriseName, defaultDatabaseName)

			clusterParams, err := expandCreateForManagedRedis(model)
			if err != nil {
				return err
			}

			if err := clusterClient.CreateCallbackThenPoll(ctx, clusterId, clusterParams, metadata.SetIDCallback(&clusterId)); err != nil {
				return fmt.Errorf("creating %s: %+v", clusterId, err)
			}

			metadata.SetID(clusterId)

			pollerType := custompollers.NewClusterStatePoller(clusterClient, clusterId)
			poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for `resourceState` to be `Running` for %s: %+v", clusterId, err)
			}

			if len(model.DefaultDatabase) == 1 {
				dbModel := model.DefaultDatabase[0]

				if err := createDb(ctx, dbClient, dbId, dbModel); err != nil {
					return fmt.Errorf("creating %s: %+v", dbId, err)
				}
			}

			return nil
		},
	}
}
